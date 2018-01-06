package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"localize"
	"transport"
)

var timezone *time.Location

func init() {
	http.HandleFunc("/dialogflow", dialogflow)
	timezone, _ = time.LoadLocation("Europe/Zurich")
}

// Returns zero time if failure.
func tryParseStupidDate(raw string) time.Time {
	var result time.Time
	var err error
	if result, err = time.Parse("15:04:05", raw); err != nil {
		if result, err = time.Parse("2006-01-02T15:04:05Z", raw); err != nil {
			result = time.Time{}
		}
	} else {
		// Successfully parsed as HH:MM:SS, so we assume it's today.
		// XXX: This will be dumb around midnight, I guess. Remember to email the Dialogflow folks about how this is silly.
		now := time.Now()
		result = time.Date(now.Year(), now.Month(), now.Day(), result.Hour(), result.Minute(), result.Second(), 0, now.Location())
	}
	if !result.IsZero() { // Make sure it's Swiss time or something.
		result = result.In(timezone)
	}
	return result
}

func mode(category string) string {
	if r, ok := map[string]string{
		"strain":        "train",
		"express_train": "train",
		"walk":          "walk",
		"tram":          "tram"}[category]; ok {
		return r
	}
	// XXX: ????
	return category
}

func prettyName(category, number string) string {
	switch category {
	case "BUS":
		return number
	case "T":
		return number
	}
	return fmt.Sprintf("%s%s", category, number)
}

func dialogflow(writer http.ResponseWriter, req *http.Request) {
	handleError := func(f string, xs ...interface{}) {
		log.Errorf(appengine.NewContext(req), f, xs...)
		http.Error(writer, fmt.Sprintf(f, xs...), http.StatusInternalServerError)
	}
	// Parse request body into DialogflowRequest
	dreq := DialogflowRequest{}
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError("Error reading POST: %v", err)
		return
	}
	if err := json.Unmarshal(bs, &dreq); err != nil {
		handleError("Error unmarshalling POST: %v", err)
		return
	}
	dresp := DialogflowResponse{}
	svc := transport.Transport{
		Client: urlfetch.Client(appengine.NewContext(req)),
		Logger: func(x string) { log.Infof(appengine.NewContext(req), "%s", x) },
	}

	switch dreq.Result.Metadata.IntentName {
	case "next-departure":
		err = stationboard(svc, dreq, &dresp)
	case "next-departures":
		err = stationboard(svc, dreq, &dresp)
	case "find-stations":
		err = findStations(svc, dreq, &dresp)
	default:
		err = fmt.Errorf("Unknown intent %s", dreq.Result.Metadata.IntentName)
	}

	if err != nil {
		handleError("%v", err)
	}
	bs, err = json.Marshal(dresp)
	if err != nil {
		handleError("Error marshalling response: %v", err)
		return
	}
	if _, err := writer.Write(bs); err != nil {
		handleError("Error writing response: %v", err)
		return
	}
}

func findStations(svc transport.Transport, dreq DialogflowRequest, dresp *DialogflowResponse) error {
	lreq := transport.LocationsRequest{
		Query: dreq.Result.Parameters.Query,
		Lat:   dreq.OriginalRequest.Data.Device.Location.Coordinates.Latitude,
		Lon:   dreq.OriginalRequest.Data.Device.Location.Coordinates.Longitude,
	}
	loc := localize.NewLocalizer(dreq.Lang, timezone)
	if lreq.Query == "" && !(lreq.Lat != 0.0 && lreq.Lon != 0.0) {
		// Request the user location.
		dresp.Speech = loc.NeedLocation()
		dresp.Data.Google.SystemIntent.Intent = "assistant.intent.action.PERMISSION"
		dresp.Data.Google.SystemIntent.Spec.PermissionValueSpec.OptContext = loc.PermissionContext()
		dresp.Data.Google.SystemIntent.Spec.PermissionValueSpec.Permissions = []string{"DEVICE_PRECISE_LOCATION"}
		return nil
	}
	lresp, err := svc.Locations(lreq)
	if err != nil {
		return fmt.Errorf("Error calling Opendata: %v", err)
	}

	stats := []localize.Station{}
	limit := 3
	if l, _ := dreq.Result.Parameters.Limit.Int64(); l > 0 {
		limit = int(l)
	}
	for _, s := range lresp {
		stats = append(stats, localize.Station{Name: s.Label, Distance: s.Dist})
		if len(stats) == limit {
			break
		}
	}
	near := dreq.Result.Parameters.Query
	if near == "" {
		// Use device location.
		near = dreq.OriginalRequest.Data.Device.Location.FormattedAddress
	}
	dresp.Speech = loc.Stations(near, stats)
	return nil
}

func stationboard(svc transport.Transport, dreq DialogflowRequest, dresp *DialogflowResponse) error {
	// XXX: Dialogflow gives us *either* 15:04:05 OR 2006-01-02T15:04:05Z. I don't know why.
	startTime := tryParseStupidDate(dreq.Result.Parameters.DateTime)

	// Fill in the departures list to localize from *either* /connections or /stationboard.
	// This lets us share the localization code.
	departures := []localize.Departure{}

	if dreq.Result.Parameters.Destination != "" {
		// Do a /connections RPC.
		creq := transport.ConnectionsRequest{
			Station:     dreq.Result.Parameters.Source,
			Destination: dreq.Result.Parameters.Destination,
			Datetime:    startTime,
		}
		cresp, err := svc.Connections(creq)
		if err != nil {
			return fmt.Errorf("Error calling Opendata: %v", err)
		}
		for _, c := range cresp.Connections {
			// Find the first non-walking departure leg.
			for _, l := range c.Legs {
				// XXX: Probably should warn people if they have to walk somewhere first.
				if l.Type == "walk" {
					continue
				}
				d := localize.Departure{
					From:     c.From,
					Name:     l.Line,
					To:       l.Exit.SbbName,
					Mode:     mode(l.Type),
					Platform: l.SbbName,
				}
				if l.DepDelay == "" {
				} else if del, err := strconv.Atoi(l.DepDelay); err != nil {
					return err
				} else {
					d.MinutesDelay = del
				}
				if tm, err := time.ParseInLocation("2006-01-02 15:04:05", l.Departure, timezone); err != nil {
					return err
				} else {
					d.Departing = tm
				}
				departures = append(departures, d)
				// Skip the following legs of the journey.
				// XXX: Probably should say SOMETHING about them.
				break
			}

		}

	} else {
		// Do a /stationboard RPC.
		sreq := transport.StationboardRequest{
			Station:  dreq.Result.Parameters.Source,
			Mode:     transport.DEPARTURE, // XXX: Hardcoded for now
			Datetime: startTime,
		}
		sresp, err := svc.Stationboard(sreq)
		if err != nil {
			return fmt.Errorf("Error calling Opendata: %v", err)
		}
		for _, c := range sresp.Connections {
			d := localize.Departure{
				From:     sresp.Stop.Name,
				Name:     c.Line,
				To:       c.Terminal.Name,
				Mode:     mode(c.Type),
				Platform: c.Track,
			}
			if c.DepDelay == "" {
			} else if del, err := strconv.Atoi(c.DepDelay); err != nil {
				return err
			} else {
				d.MinutesDelay = del
			}
			if tm, err := time.ParseInLocation("2006-01-02 15:04:05", c.Time, timezone); err != nil {
				return err
			} else {
				d.Departing = tm
			}
			departures = append(departures, d)
		}
	}
	loc := localize.NewLocalizer(dreq.Lang, timezone)
	limit := 5 // Default
	if i, err := dreq.Result.Parameters.Limit.Int64(); err == nil {
		limit = int(i)
	}

	allowedModes := map[string]bool{}
	for _, tp := range dreq.Result.Parameters.Transport {
		allowedModes[tp] = true
	}

	// Filter "departures."
	filtered := []localize.Departure{}
	for _, d := range departures {
		// If the user specified specific routes, skip on that basis.
		if len(dreq.Result.Parameters.Route) > 0 {
			ok := false
			for _, r := range dreq.Result.Parameters.Route {
				if d.Name == r {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}
		// Or if the user specified modes.
		if len(allowedModes) > 0 {
			if !allowedModes[d.Mode] {
				continue
			}
		}
		filtered = append(filtered, d)
		if len(filtered) == limit {
			break
		}
	}
	dresp.Speech = loc.NextDepartures(dreq.Result.Parameters.Source, startTime, filtered)
	return nil
}
