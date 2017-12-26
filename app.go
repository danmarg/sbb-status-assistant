package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/dialogflow", dialogflow)
}

func mode(s StationboardStation) string {
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

func prettyName(s StationboardStation) string {
	switch s.Category {
	case "BUS":
		return s.Number
	case "T":
		return s.Number
	}
	return fmt.Sprintf("%s%s", s.Category, s.Number)
}

func userGivenName(s StationboardStation) string {
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
	svc := Transport{
		Client: urlfetch.Client(appengine.NewContext(req)),
		Logger: func(x string) { log.Infof(appengine.NewContext(req), x, nil) },
	}
	sreq := StationboardRequest{
		Station: dreq.Result.Parameters.ZvvStops,
		Type:    DEPARTURE, // XXX: Hardcoded for now
	}
	sresp, err := svc.Stationboard(sreq)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error calling Opendata: %v", err), http.StatusInternalServerError)
		return
	}
	loc := newLocalizer(dreq.Lang)
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
	departures := []departure{}
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
		d := departure{
			Name:   prettyName(c),
			OnTime: c.Stop.Prognosis.Departure != "",
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

	dresp.Speech = loc.nextDepartures(dreq.Result.Parameters.ZvvStops, departures)
}
