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

func mode(s transport.StationboardStation) string {
	switch s.Category {
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
	return s.Category
}

func prettyName(s transport.StationboardStation) string {
	switch s.Category {
	case "BUS":
		return s.Number
	case "T":
		return s.Number
	}
	return fmt.Sprintf("%s%s", s.Category, s.Number)
}

func userGivenName(s transport.StationboardStation) string {
	switch s.Category {
	case "T":
		return fmt.Sprintf("%s", s.Number)
	case "BUS":
		return fmt.Sprintf("%s", s.Number)
	default:
		return fmt.Sprintf("%s%s", s.Category, s.Number)
	}
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
	// Then dispatch to Opendata
	svc := transport.Transport{
		Client: urlfetch.Client(appengine.NewContext(req)),
		Logger: func(x string) { log.Infof(appengine.NewContext(req), x) },
	}
	var startTime time.Time
	// XXX: Dialogflow gives us *either* 15:04:05 OR 2006-01-02T15:04:05Z. I don't know why.
	startTime = tryParseStupidDate(dreq.Result.Parameters.DateTime)
	sreq := transport.StationboardRequest{
		Station:  dreq.Result.Parameters.ZvvStops,
		Type:     transport.DEPARTURE, // XXX: Hardcoded for now
		Datetime: startTime,
	}
	sresp, err := svc.Stationboard(sreq)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error calling Opendata: %v", err), http.StatusInternalServerError)
		return
	}
	loc := localize.NewLocalizer(dreq.Lang)
	// Then create response
	dresp := DialogflowResponse{}
	defer func() {
		bs, err := json.Marshal(dresp)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
			return
		}
		if _, err := writer.Write(bs); err != nil {
			http.Error(writer, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
			return
		}
	}()
	limit := 5 // Default
	if i, err := dreq.Result.Parameters.Cardinal.Int64(); err == nil {
		limit = int(i)
	}

	allowedModes := map[string]bool{}
	for _, tp := range dreq.Result.Parameters.Transport {
		allowedModes[tp] = true
	}
	departures := []localize.Departure{}
	for _, c := range sresp.Stationboard {
		// If the user specified specific routes, skip on that basis.
		if len(dreq.Result.Parameters.ZvvRoutes) > 0 {
			ok := false
			n := userGivenName(c)
			for _, r := range dreq.Result.Parameters.ZvvRoutes {
				if n == r {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}
		// Or if the user specified modes.
		if len(allowedModes) > 0 {
			if !allowedModes[mode(c)] {
				continue
			}
		}
		d := localize.Departure{
			Name:   prettyName(c),
			OnTime: c.Stop.Prognosis.Departure == "",
			To:     c.To,
			Mode:   mode(c),
		}
		if dp, err := time.Parse("2006-01-02T15:04:05-07:00", c.Stop.Prognosis.Departure); c.Stop.Prognosis.Departure != "" && err != nil {
			d.Departing = dp
		} else {
			d.Departing = time.Unix(int64(c.Stop.DepartureTimestamp), 0)
		}
		departures = append(departures, d)
		if len(departures) == limit {
			break
		}
	}

	dresp.Speech = loc.NextDepartures(dreq.Result.Parameters.ZvvStops, startTime, departures)
}
