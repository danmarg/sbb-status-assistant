package transport

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const stationboardEndpoint = "http://transport.opendata.ch/v1/stationboard"

type Transport struct {
	Client *http.Client
	Logger func(string)
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

	strParams := []string{}
	for k, v := range params {
		strParams = append(strParams, k+"="+url.QueryEscape(v))
	}

	u := stationboardEndpoint + "?" + strings.Join(strParams, "&")
	if t.Logger != nil {
		t.Logger("OpenTransport URL: " + u)
	}
	rq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return StationboardResponse{}, err
	}

	rsp, err := t.Client.Do(rq)
	if err != nil {
		return StationboardResponse{}, err
	}

	defer rsp.Body.Close()

	var resp StationboardResponse
	if err := json.NewDecoder(rsp.Body).Decode(&resp); err != nil {
		return StationboardResponse{}, err
	}

	// Post-filter by routes, since this isn't supported by the Opendata.ch API.
	// XXX: I think we can filter by the Number field...?
	return resp, nil
}
