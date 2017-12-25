package transport

import (
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/dialogflow", dialogflow)
}

func dialogflow(writer http.ResponseWriter, req *http.Request) {
	// Parse request body into DialogflowRequest
	dreq := DialogflowRequest{}
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(bs, &dreq); err != nil {
		log.Fatalln(err)
	}
	// Then dispatch to Opendata
	svc := Transport{
		Client: urlfetch.Client(appengine.NewContext(req)),
	}
	sreq := StationboardRequest{
		Station: dreq.Result.Parameters.ZvvStops,
		Limit:   dreq.Result.Parameters.Cardinal,
		Type:    DEPARTURES, // XXX: Hardcoded for now
	}
	tps := make([]string, len(dreq.Result.Parameters.Transport))
	for i, tp := range dreq.Result.Parameters.Transport {
		switch tp {
		case "tram":
			tps[i] = TRAM
		case "bus":
			tps[i] = BUS
		case "train":
			tps[i] = TRAIN
		case "boat":
			tps[i] = BOAT
		case "any":
			tps[i] = ANY
		}
	}
	sreq.Transportations = tps
	sresp, err := svc.Stationboard(sreq)
	if err != nil {
		log.Fatal(err)
	}
	// Then create response
}
