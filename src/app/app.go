package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"localize"
	"transport"
)

func init() {
	http.HandleFunc("/dialogflow", dialogflow)
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
		l, _ := time.LoadLocation("Europe/Zurich")
		result = result.In(l)
	}
	return result
}

func mode(category string) string {
	switch category {
	case "BUS":
		return "bus"
	case "T":
		return "tram"
	case "IC":
	case "IR":
	case "S":
		return "train"
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
	// Parse request body into DialogflowRequest
	dreq := DialogflowRequest{}
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error reading POST: %v", err), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(bs, &dreq); err != nil {
		http.Error(writer, fmt.Sprintf("Error unmarshalling POST: %v", err), http.StatusInternalServerError)
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
	default:
		err = fmt.Errorf("Unknown intent %s", dreq.Result.Metadata.IntentName)
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	bs, err = json.Marshal(dresp)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := writer.Write(bs); err != nil {
		http.Error(writer, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}
}

func findStations(svc transport.Transport, dreq DialogflowRequest, dresp *DialogflowResponse) error {
	lreq := transport.LocationsRequest{
		Query: dreq.Result.Parameters.Query,
		Lat:   dreq.OriginalRequest.Data.Device.Location.Coordinates.Latitude,
		Lon:   dreq.OriginalRequest.Data.Device.Location.Coordinates.Longitude,
	}
	loc := localize.NewLocalizer(dreq.Lang)
	if lreq.Query == "" && !(lreq.Lat != 0.0 && lreq.Lon != 0.0) {
		// Request the user location.
		dresp.Speech = loc.NeedLocation()
		dresp.Data.Google.PermissionsRequest.Permissions = []string{"DEVICE_PRECISE_LOCATION"}
		return nil
	}
	lresp, err := svc.Locations(lreq)
	if err != nil {
		return fmt.Errorf("Error calling Opendata: %v", err)
	}

	stats := []localize.Station{}
	for _, s := range lresp.Stations {
		stats = append(stats, localize.Station{Name: s.Name, Distance: s.Distance})
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
			// XXX: Probably should support multiple-connection paths at some point.
			nonWalking := 0
			for _, j := range c.Sections {
				if j.Walk.Duration > 0 {
					nonWalking++
				}
			}
			if nonWalking > 1 {
				continue
			}
			d := localize.Departure{
				Name:     prettyName(c.Sections[0].Journey.Category, c.Sections[0].Journey.Number),
				To:       c.To.Station.Name,
				OnTime:   true,
				Mode:     mode(c.Sections[0].Journey.Category),
				Platform: c.From.Platform,
			}
			if x, _ := c.From.Delay.Int64(); c.From.Delay.String() != "" && x > 0 {
				d.OnTime = false
			}
			if dp, err := time.Parse("2006-01-02T15:04:05-07:00", c.From.Prognosis.Departure); c.From.Prognosis.Departure != "" && err != nil {
				d.Departing = dp
			} else {
				d.Departing = time.Unix(int64(c.From.DepartureTimestamp), 0)
			}
			departures = append(departures, d)
		}

	} else {
		// Do a /stationboard RPC.
		sreq := transport.StationboardRequest{
			Station:  dreq.Result.Parameters.Source,
			Type:     transport.DEPARTURE, // XXX: Hardcoded for now
			Datetime: startTime,
		}
		sresp, err := svc.Stationboard(sreq)
		if err != nil {
			return fmt.Errorf("Error calling Opendata: %v", err)
		}
		for _, c := range sresp.Stationboard {
			d := localize.Departure{
				Name:     prettyName(c.Category, c.Number),
				OnTime:   c.Stop.Prognosis.Departure == "",
				To:       c.To,
				Mode:     mode(c.Category),
				Platform: c.Stop.Platform,
			}
			if dp, err := time.Parse("2006-01-02T15:04:05-07:00", c.Stop.Prognosis.Departure); c.Stop.Prognosis.Departure != "" && err != nil {
				d.Departing = dp
			} else {
				d.Departing = time.Unix(int64(c.Stop.DepartureTimestamp), 0)
			}
			departures = append(departures, d)
		}

	}
	loc := localize.NewLocalizer(dreq.Lang)
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
