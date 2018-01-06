package app

import (
	"encoding/json"
	"io/ioutil"
)

const (
	stationEntitiesFile = "static/station_entities.json"
)

type entityMap map[string]string

func loadEntities() (entityMap, error) {
	var r entityMap
	if bs, err := ioutil.ReadFile(stationEntitiesFile); err != nil {
		return r, err
	} else {
		err := json.Unmarshal(bs, &r)
		return r, err
	}
}
