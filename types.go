package transport

import (
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
					Type string `json:"type"`
					X    string `json:"x"`
					Y    string `json:"y"`
				} `json:"coordinate"`
			} `json:"station"`
			Arrival            interface{} `json:"arrival"`
			ArrivalTimestamp   interface{} `json:"arrivalTimestamp"`
			Departure          time.Time   `json:"departure"`
			DepartureTimestamp int         `json:"departureTimestamp"`
			Platform           string      `json:"platform"`
			Prognosis          struct {
				Platform    interface{} `json:"platform"`
				Arrival     interface{} `json:"arrival"`
				Departure   interface{} `json:"departure"`
				Capacity1St string      `json:"capacity1st"`
				Capacity2Nd string      `json:"capacity2nd"`
			} `json:"prognosis"`
		} `json:"stop"`
		Name     string      `json:"name"`
		Category string      `json:"category"`
		Number   string      `json:"number"`
		Operator interface{} `json:"operator"`
		To       string      `json:"to"`
	} `json:"stationboard"`
}

type DialogflowRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Lang      string    `json:"lang"`
	Result    struct {
		Source           string `json:"source"`
		ResolvedQuery    string `json:"resolvedQuery"`
		Action           string `json:"action"`
		ActionIncomplete bool   `json:"actionIncomplete"`
		Parameters       struct {
			ZvvStops  []string      `json:"zvv_stops"`
			Transport []string      `json:"transport"`
			ZvvRoutes []interface{} `json:"zvv_routes"`
			Cardinal  int           `json:"cardinal"`
		} `json:"parameters"`
		Contexts []interface{} `json:"contexts"`
		Metadata struct {
			IntentID                  string `json:"intentId"`
			WebhookUsed               string `json:"webhookUsed"`
			WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
			WebhookResponseTime       int    `json:"webhookResponseTime"`
			IntentName                string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech   string `json:"speech"`
			Messages []struct {
				Type   int    `json:"type"`
				Speech string `json:"speech"`
			} `json:"messages"`
		} `json:"fulfillment"`
		Score int `json:"score"`
	} `json:"result"`
	Status struct {
		Code            int    `json:"code"`
		ErrorType       string `json:"errorType"`
		ErrorDetails    string `json:"errorDetails"`
		WebhookTimedOut bool   `json:"webhookTimedOut"`
	} `json:"status"`
	SessionID string `json:"sessionId"`
}

type DialogflowResponse struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Source      string `json:"source"`
}
