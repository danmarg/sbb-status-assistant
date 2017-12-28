package transport

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	stationboardEndpoint = "http://transport.opendata.ch/v1/stationboard"
	connectionsEndpoint  = "http://transport.opendata.ch/v1/connections"
	locationsEndpoint    = "http://transport.opendata.ch/v1/locations"
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
		params["query"] = req.Query
	}
	if req.Lat != 0.0 {
		params["x"] = strconv.FormatFloat(req.Lat, 'f', -1, 32)
	}
	if req.Lon != 0.0 {
		params["y"] = strconv.FormatFloat(req.Lon, 'f', -1, 32)
	}
	for i, tp := range req.Transportations {
		if i > 0 {
			params["transportations"] += ","
		}
		params["transportations"] += transportations[tp]
	}

	var resp LocationsResponse
	err := t.dispatch(locationsEndpoint, params, &resp)
	return resp, err
}

func (t *Transport) Stationboard(req StationboardRequest) (StationboardResponse, error) {
	params := map[string]string{}
	if req.Station != "" {
		params["station"] = req.Station
	}
	if req.Id != "" {
		params["id"] = req.Id
	}
	if req.Limit != 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	for i, tp := range req.Transportations {
		if i > 0 {
			params["transportations"] += ","
		}
		params["transportations"] += transportations[tp]
	}
	if !req.Datetime.IsZero() {
		params["datetime"] = req.Datetime.Format("2006-01-02 15:04")
	}
	if req.Type != 0 {
		params["type"] = types[req.Type]
	}
	var resp StationboardResponse
	err := t.dispatch(stationboardEndpoint, params, &resp)
	// Post-filter by routes, since this isn't supported by the Opendata.ch API.
	// XXX: I think we can filter by the Number field...?
	return resp, err
}

func (t *Transport) Connections(req ConnectionsRequest) (ConnectionsResponse, error) {
	params := map[string]string{}
	if req.Station != "" {
		params["from"] = req.Station
	}
	if req.Station != "" {
		params["to"] = req.Destination
	}
	if req.Limit != 0 {
		params["limit"] = strconv.Itoa(req.Limit)
	}
	for i, tp := range req.Transportations {
		if i > 0 {
			params["transportations"] += ","
		}
		params["transportations"] += transportations[tp]
	}
	if !req.Datetime.IsZero() {
		params["date"] = req.Datetime.Format("2006-01-02")
		params["time"] = req.Datetime.Format("15:04")
	}
	var resp ConnectionsResponse
	err := t.dispatch(connectionsEndpoint, params, &resp)

	// Post-filter by routes, since this isn't supported by the Opendata.ch API.
	// XXX: I think we can filter by the Number field...?
	return resp, err
}
