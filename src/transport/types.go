package transport

import (
	"encoding/json"
	"time"
)

const (
	_ = iota
	// Transportations:
	TRAM  = iota
	BUS   = iota
	BOAT  = iota
	TRAIN = iota
	ANY   = iota
	// Types:
	DEPARTURE = iota
	ARRIVAL   = iota
)

var (
	transportations = map[int]string{
		TRAM:  "tramway_underground",
		BUS:   "bus",
		BOAT:  "ship",
		TRAIN: "ice_tgv_rj,ec_ic,ir,re_d,s_sn_r,arz_ext",
		ANY:   "",
	}

	types = map[int]string{
		DEPARTURE: "departure",
		ARRIVAL:   "arrival",
	}
)

type StationboardRequest struct {
	Station         string
	Id              string
	Limit           int
	Transportations []int // One of the consts above
	Datetime        time.Time
	Type            int // One of the consts above
	// Optional route
	Route string
}

type StationboardStation struct {
	Stop struct {
		Station struct {
			ID         string      `json:"id"`
			Name       string      `json:"name"`
			Score      interface{} `json:"score"`
			Coordinate struct {
				Type string      `json:"type"`
				X    json.Number `json:"x"`
				Y    json.Number `json:"y"`
			} `json:"coordinate"`
		} `json:"station"`
		Arrival            string `json:"arrival"`
		ArrivalTimestamp   int    `json:"arrivalTimestamp"`
		Departure          string `json:"departure"`
		DepartureTimestamp int    `json:"departureTimestamp"`
		Platform           string `json:"platform"`
		Prognosis          struct {
			Platform    string `json:"platform"`
			Arrival     string `json:"arrival"`
			Departure   string `json:"departure"`
			Capacity1St string `json:"capacity1st"`
			Capacity2Nd string `json:"capacity2nd"`
		} `json:"prognosis"`
	} `json:"stop"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Number   string `json:"number"`
	Operator string `json:"operator"`
	To       string `json:"to"`
}

type StationboardResponse struct {
	Stationboard []StationboardStation `json:"stationboard"`
}
