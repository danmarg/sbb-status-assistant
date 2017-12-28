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

type StationboardResponse struct {
	Stationboard []struct {
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
	} `json:"stationboard"`
}

type LocationsRequest struct {
	Query           string
	Lat             float64
	Lon             float64
	Transportations []int // One of the consts above
	// Note that the API also takes a "Type" parameter, but here we only want to support station lookups.
}

type LocationsResponse struct {
	Stations []struct {
		ID         interface{} `json:"id"`
		Name       string      `json:"name"`
		Score      interface{} `json:"score"`
		Coordinate struct {
			Type string  `json:"type"`
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
		} `json:"coordinate"`
		Distance int `json:"distance"`
	} `json:"stations"`
}

type ConnectionsRequest struct {
	Station         string
	Destination     string
	Limit           int
	Transportations []int // One of the consts above
	Datetime        time.Time
	// Optional route
	Route string
}

type ConnectionsResponse struct {
	Connections []struct {
		From struct {
			Station struct {
				ID         string      `json:"id"`
				Name       string      `json:"name"`
				Score      interface{} `json:"score"`
				Coordinate struct {
					Type string  `json:"type"`
					X    float64 `json:"x"`
					Y    float64 `json:"y"`
				} `json:"coordinate"`
				Distance interface{} `json:"distance"`
			} `json:"station"`
			Arrival            interface{} `json:"arrival"`
			ArrivalTimestamp   interface{} `json:"arrivalTimestamp"`
			Departure          string      `json:"departure"`
			DepartureTimestamp int         `json:"departureTimestamp"`
			Delay              json.Number `json:"delay"`
			Platform           string      `json:"platform"`
			Prognosis          struct {
				Platform    interface{} `json:"platform"`
				Arrival     interface{} `json:"arrival"`
				Departure   string      `json:"departure"`
				Capacity1St interface{} `json:"capacity1st"`
				Capacity2Nd interface{} `json:"capacity2nd"`
			} `json:"prognosis"`
			RealtimeAvailability interface{} `json:"realtimeAvailability"`
			Location             struct {
				ID         string      `json:"id"`
				Name       string      `json:"name"`
				Score      interface{} `json:"score"`
				Coordinate struct {
					Type string  `json:"type"`
					X    float64 `json:"x"`
					Y    float64 `json:"y"`
				} `json:"coordinate"`
				Distance interface{} `json:"distance"`
			} `json:"location"`
		} `json:"from"`
		To struct {
			Station struct {
				ID         string      `json:"id"`
				Name       string      `json:"name"`
				Score      interface{} `json:"score"`
				Coordinate struct {
					Type string  `json:"type"`
					X    float64 `json:"x"`
					Y    float64 `json:"y"`
				} `json:"coordinate"`
				Distance interface{} `json:"distance"`
			} `json:"station"`
			Arrival            string      `json:"arrival"`
			ArrivalTimestamp   int         `json:"arrivalTimestamp"`
			Departure          interface{} `json:"departure"`
			DepartureTimestamp interface{} `json:"departureTimestamp"`
			Delay              json.Number `json:"delay"`
			Platform           interface{} `json:"platform"`
			Prognosis          struct {
				Platform    interface{} `json:"platform"`
				Arrival     interface{} `json:"arrival"`
				Departure   interface{} `json:"departure"`
				Capacity1St interface{} `json:"capacity1st"`
				Capacity2Nd interface{} `json:"capacity2nd"`
			} `json:"prognosis"`
			RealtimeAvailability interface{} `json:"realtimeAvailability"`
			Location             struct {
				ID         string      `json:"id"`
				Name       string      `json:"name"`
				Score      interface{} `json:"score"`
				Coordinate struct {
					Type string  `json:"type"`
					X    float64 `json:"x"`
					Y    float64 `json:"y"`
				} `json:"coordinate"`
				Distance interface{} `json:"distance"`
			} `json:"location"`
		} `json:"to"`
		Duration    string      `json:"duration"`
		Transfers   int         `json:"transfers"`
		Service     interface{} `json:"service"`
		Products    []string    `json:"products"`
		Capacity1St interface{} `json:"capacity1st"`
		Capacity2Nd interface{} `json:"capacity2nd"`
		Sections    []struct {
			Journey struct {
				Name         string      `json:"name"`
				Category     string      `json:"category"`
				Subcategory  interface{} `json:"subcategory"`
				CategoryCode interface{} `json:"categoryCode"`
				Number       string      `json:"number"`
				Operator     string      `json:"operator"`
				To           string      `json:"to"`
				PassList     []struct {
					Station struct {
						ID         string      `json:"id"`
						Name       string      `json:"name"`
						Score      interface{} `json:"score"`
						Coordinate struct {
							Type string  `json:"type"`
							X    float64 `json:"x"`
							Y    float64 `json:"y"`
						} `json:"coordinate"`
						Distance interface{} `json:"distance"`
					} `json:"station"`
					Arrival            interface{} `json:"arrival"`
					ArrivalTimestamp   interface{} `json:"arrivalTimestamp"`
					Departure          string      `json:"departure"`
					DepartureTimestamp int         `json:"departureTimestamp"`
					Delay              json.Number `json:"delay"`
					Platform           string      `json:"platform"`
					Prognosis          struct {
						Platform    interface{} `json:"platform"`
						Arrival     interface{} `json:"arrival"`
						Departure   string      `json:"departure"`
						Capacity1St interface{} `json:"capacity1st"`
						Capacity2Nd interface{} `json:"capacity2nd"`
					} `json:"prognosis"`
					RealtimeAvailability interface{} `json:"realtimeAvailability"`
					Location             struct {
						ID         string      `json:"id"`
						Name       string      `json:"name"`
						Score      interface{} `json:"score"`
						Coordinate struct {
							Type string  `json:"type"`
							X    float64 `json:"x"`
							Y    float64 `json:"y"`
						} `json:"coordinate"`
						Distance interface{} `json:"distance"`
					} `json:"location"`
				} `json:"passList"`
				Capacity1St interface{} `json:"capacity1st"`
				Capacity2Nd interface{} `json:"capacity2nd"`
			} `json:"journey"`
			Walk struct {
				Duration int `json:"duration"`
			} `json:"walk"`
			Departure struct {
				Station struct {
					ID         string      `json:"id"`
					Name       string      `json:"name"`
					Score      interface{} `json:"score"`
					Coordinate struct {
						Type string  `json:"type"`
						X    float64 `json:"x"`
						Y    float64 `json:"y"`
					} `json:"coordinate"`
					Distance interface{} `json:"distance"`
				} `json:"station"`
				Arrival            interface{} `json:"arrival"`
				ArrivalTimestamp   interface{} `json:"arrivalTimestamp"`
				Departure          string      `json:"departure"`
				DepartureTimestamp int         `json:"departureTimestamp"`
				Delay              json.Number `json:"delay"`
				Platform           string      `json:"platform"`
				Prognosis          struct {
					Platform    interface{} `json:"platform"`
					Arrival     interface{} `json:"arrival"`
					Departure   string      `json:"departure"`
					Capacity1St interface{} `json:"capacity1st"`
					Capacity2Nd interface{} `json:"capacity2nd"`
				} `json:"prognosis"`
				RealtimeAvailability interface{} `json:"realtimeAvailability"`
				Location             struct {
					ID         string      `json:"id"`
					Name       string      `json:"name"`
					Score      interface{} `json:"score"`
					Coordinate struct {
						Type string  `json:"type"`
						X    float64 `json:"x"`
						Y    float64 `json:"y"`
					} `json:"coordinate"`
					Distance interface{} `json:"distance"`
				} `json:"location"`
			} `json:"departure"`
			Arrival struct {
				Station struct {
					ID         string      `json:"id"`
					Name       string      `json:"name"`
					Score      interface{} `json:"score"`
					Coordinate struct {
						Type string  `json:"type"`
						X    float64 `json:"x"`
						Y    float64 `json:"y"`
					} `json:"coordinate"`
					Distance interface{} `json:"distance"`
				} `json:"station"`
				Arrival            string      `json:"arrival"`
				ArrivalTimestamp   int         `json:"arrivalTimestamp"`
				Departure          interface{} `json:"departure"`
				DepartureTimestamp interface{} `json:"departureTimestamp"`
				Delay              json.Number `json:"delay"`
				Platform           string      `json:"platform"`
				Prognosis          struct {
					Platform    interface{} `json:"platform"`
					Arrival     string      `json:"arrival"`
					Departure   interface{} `json:"departure"`
					Capacity1St interface{} `json:"capacity1st"`
					Capacity2Nd interface{} `json:"capacity2nd"`
				} `json:"prognosis"`
				RealtimeAvailability interface{} `json:"realtimeAvailability"`
				Location             struct {
					ID         string      `json:"id"`
					Name       string      `json:"name"`
					Score      interface{} `json:"score"`
					Coordinate struct {
						Type string  `json:"type"`
						X    float64 `json:"x"`
						Y    float64 `json:"y"`
					} `json:"coordinate"`
					Distance interface{} `json:"distance"`
				} `json:"location"`
			} `json:"arrival"`
		} `json:"sections"`
	} `json:"connections"`
	From struct {
		ID         string      `json:"id"`
		Name       string      `json:"name"`
		Score      interface{} `json:"score"`
		Coordinate struct {
			Type string  `json:"type"`
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
		} `json:"coordinate"`
		Distance interface{} `json:"distance"`
	} `json:"from"`
	To struct {
		ID         string      `json:"id"`
		Name       string      `json:"name"`
		Score      interface{} `json:"score"`
		Coordinate struct {
			Type string  `json:"type"`
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
		} `json:"coordinate"`
		Distance interface{} `json:"distance"`
	} `json:"to"`
	Stations struct {
		From []struct {
			ID         string      `json:"id"`
			Name       string      `json:"name"`
			Score      interface{} `json:"score"`
			Coordinate struct {
				Type string  `json:"type"`
				X    float64 `json:"x"`
				Y    float64 `json:"y"`
			} `json:"coordinate"`
			Distance interface{} `json:"distance"`
		} `json:"from"`
		To []struct {
			ID         string      `json:"id"`
			Name       string      `json:"name"`
			Score      interface{} `json:"score"`
			Coordinate struct {
				Type string  `json:"type"`
				X    float64 `json:"x"`
				Y    float64 `json:"y"`
			} `json:"coordinate"`
			Distance interface{} `json:"distance"`
		} `json:"to"`
	} `json:"stations"`
}
