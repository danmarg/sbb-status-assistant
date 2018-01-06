package transport

import (
	"encoding/json"
	"time"
)

const (
	_         = iota
	ARRIVAL   = iota
	DEPARTURE = iota
)

type StationboardRequest struct {
	Station  string
	Limit    int
	Datetime time.Time
	Mode     int // ARRIVAL or DEPARTURE
}

type StationboardResponse struct {
	Stop struct {
		ID   string      `json:"id"`
		Name string      `json:"name"`
		X    json.Number `json:"x"`
		Y    json.Number `json:"y"`
	} `json:"stop"`
	Connections []struct {
		Time     string `json:"time"`
		G        string `json:"*G"`
		L        string `json:"*L"`
		Type     string `json:"type"`
		Line     string `json:"line"`
		Operator string `json:"operator"`
		Color    string `json:"color"`
		Number   string `json:"number"`
		Terminal struct {
			ID   string      `json:"id"`
			Name string      `json:"name"`
			X    json.Number `json:"x"`
			Y    json.Number `json:"y"`
		} `json:"terminal"`
		SubsequentStops []struct {
			ID  string      `json:"id"`
			X   json.Number `json:"x"`
			Y   json.Number `json:"y"`
			Arr string      `json:"arr"`
			Dep string      `json:"dep,omitempty"`
		} `json:"subsequent_stops"`
		Track    string `json:"track,omitempty"`
		ArrDelay string `json:"arr_delay,omitempty"`
		DepDelay string `json:"dep_delay,omitempty"`
	} `json:"connections"`
	Request string `json:"request"`
	EOF     int    `json:"eof"`
}

type LocationsRequest struct {
	Query string
	Lat   float64
	Lon   float64
}

type LocationsResponse []struct {
	Label     string  `json:"label"`
	Dist      float64 `json:"dist"`
	Iconclass string  `json:"iconclass"`
}

type ConnectionsRequest struct {
	Station     string
	Destination string
	Via         string
	Limit       int
	Datetime    time.Time
}

type ConnectionsResponse struct {
	Count       int `json:"count"`
	Rawtime     int `json:"rawtime"`
	Maxtime     int `json:"maxtime"`
	Connections []struct {
		From      string `json:"from"`
		Departure string `json:"departure"`
		DepDelay  string `json:"dep_delay,omitempty"`
		To        string `json:"to"`
		Arrival   string `json:"arrival"`
		Duration  int    `json:"duration"`
		Legs      []struct {
			Departure string      `json:"departure,omitempty"`
			Tripid    string      `json:"tripid,omitempty"`
			Number    string      `json:"number,omitempty"`
			Stopid    string      `json:"stopid,omitempty"`
			X         json.Number `json:"x,omitempty"`
			Y         json.Number `json:"y,omitempty"`
			Name      string      `json:"name"`
			SbbName   string      `json:"sbb_name,omitempty"`
			Type      string      `json:"type,omitempty"`
			Line      string      `json:"line,omitempty"`
			Terminal  string      `json:"terminal,omitempty"`
			Fgcolor   string      `json:"fgcolor,omitempty"`
			Bgcolor   string      `json:"bgcolor,omitempty"`
			G         string      `json:"*G,omitempty"`
			L         string      `json:"*L,omitempty"`
			Operator  string      `json:"operator,omitempty"`
			Stops     []struct {
				Arrival   string      `json:"arrival"`
				Departure string      `json:"departure"`
				DepDelay  string      `json:"dep_delay"`
				Name      string      `json:"name"`
				Stopid    string      `json:"stopid"`
				X         json.Number `json:"x"`
				Y         json.Number `json:"y"`
			} `json:"stops,omitempty"`
			Runningtime int `json:"runningtime,omitempty"`
			Exit        struct {
				Arrival  string      `json:"arrival"`
				Stopid   string      `json:"stopid"`
				X        json.Number `json:"x"`
				Y        json.Number `json:"y"`
				Name     string      `json:"name"`
				SbbName  string      `json:"sbb_name"`
				Waittime int         `json:"waittime"`
				Track    string      `json:"track"`
				ArrDelay string      `json:"arr_delay"`
			} `json:"exit,omitempty"`
			DepDelay   string `json:"dep_delay,omitempty"`
			Track      string `json:"track,omitempty"`
			Arrival    string `json:"arrival,omitempty"`
			Waittime   int    `json:"waittime,omitempty"`
			NormalTime int    `json:"normal_time,omitempty"`
			Isaddress  bool   `json:"isaddress,omitempty"`
		} `json:"legs"`
	} `json:"connections"`
	URL    string `json:"url"`
	Points []struct {
		Text string      `json:"text"`
		URL  string      `json:"url"`
		ID   string      `json:"id,omitempty"`
		X    json.Number `json:"x,omitempty"`
		Y    json.Number `json:"y,omitempty"`
	} `json:"points"`
	Description string `json:"description"`
	Request     string `json:"request"`
	EOF         int    `json:"eof"`
}
