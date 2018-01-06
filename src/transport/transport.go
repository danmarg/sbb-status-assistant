package transport

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	stationboardEndpoint = "https://timetable.search.ch/api/stationboard.json"
	connectionsEndpoint  = "https://timetable.search.ch/api/route.json"
	locationsEndpoint    = "https://timetable.search.ch/api/completion.json"
)

type Transport struct {
	Client *http.Client
	Logger func(string)
}

func (t *Transport) dispatch(endpoint string, params map[string]string, result interface{}) error {
	strParams := []string{}
	for k, v := range params {
		strParams = append(strParams, k+"="+url.QueryEscape(v))
	}

	u := endpoint + "?" + strings.Join(strParams, "&")
	if t.Logger != nil {
		t.Logger("OpenTransport URL: " + u)
	}
	rq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	rsp, err := t.Client.Do(rq)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	return json.NewDecoder(rsp.Body).Decode(result)
}

func (t *Transport) Locations(req LocationsRequest) (LocationsResponse, error) {
	params := map[string]string{}
	if req.Query != "" {
		params["term"] = req.Query
	}
	if req.Lat != 0.0 && req.Lon != 0.0 {
		params["latlon"] = strconv.FormatFloat(req.Lat, 'f', -1, 32) + "," +
			strconv.FormatFloat(req.Lon, 'f', -1, 32)
	}

	var resp LocationsResponse
	err := t.dispatch(locationsEndpoint, params, &resp)
	return resp, err
}

func (t *Transport) Stationboard(req StationboardRequest) (StationboardResponse, error) {
	params := map[string]string{}
	if req.Station != "" {
		params["stop"] = req.Station
	}
	if req.Limit != 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	if !req.Datetime.IsZero() {
		params["date"] = req.Datetime.Format("2006-01-02")
		params["time"] = req.Datetime.Format("15:04")
	}
	if req.Mode == ARRIVAL {
		params["mode"] = "arrival"
	} else {
		params["mode"] = "depart"
	}
	params["show_tracks"] = "true"
	params["show_delays"] = "true"
	params["show_subsequent_stops"] = "true"
	params["show_trackchanges"] = "true"

	var resp StationboardResponse
	err := t.dispatch(stationboardEndpoint, params, &resp)
	return resp, err
}

func (t *Transport) Connections(req ConnectionsRequest) (ConnectionsResponse, error) {
	params := map[string]string{}
	if req.Station != "" {
		params["from"] = req.Station
	}
	if req.Destination != "" {
		params["to"] = req.Destination
	}
	if req.Via != "" {
		params["via"] = req.Via
	}
	if req.Limit != 0 {
		params["num"] = strconv.Itoa(req.Limit)
	}
	if !req.Datetime.IsZero() {
		params["date"] = req.Datetime.Format("2006-01-02")
		params["time"] = req.Datetime.Format("15:04")
	}
	params["show_delays"] = "true"
	params["show_trackchanges"] = "true"

	var resp ConnectionsResponse
	err := t.dispatch(connectionsEndpoint, params, &resp)
	return resp, err
}
