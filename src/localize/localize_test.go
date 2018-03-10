package localize

import (
	"testing"
	"time"
)

func TestNeedLocation(t *testing.T) {
	for l, want := range map[string]string{
		"en": "I need your location.",
		"de": "Ich brauche Ihren Standort.",
	} {
		l := NewLocalizer(l, time.Now().Location())
		got := l.NeedLocation()
		if got != want {
			t.Errorf("want '%v', got '%v'", want, got)
		}
	}
}

func TestPermissionContext(t *testing.T) {
	for l, want := range map[string]string{
		"en": "To look for stations",
		"de": "Um Haltestellen zu suchen",
	} {
		l := NewLocalizer(l, time.Now().Location())
		got := l.PermissionContext()
		if got != want {
			t.Errorf("want '%v', got '%v'", want, got)
		}
	}
}

type stationsTest struct {
	Near     string
	Stations []Station
	Want     string
}

func TestStations(t *testing.T) {
	for l, wants := range map[string][]stationsTest{
		"en": []stationsTest{
			{"", []Station{}, "I could not find any matching stations."},
			{"Zurich", []Station{}, "I could not find any matching stations near Zurich."},
			{"", []Station{{"Zurich HB", 1.2}}, "The closest station to you is: Zurich HB, 1 meter away."},
			{"Zurich", []Station{{"Zurich HB", 1.2}}, "The closest station to Zurich is: Zurich HB, 1 meter away."},
			{"", []Station{{"Zurich HB SZU", 2.7}, {"Zurich HB", 1.2}}, "The closest stations to you are: Zurich HB SZU, 2 meters away; Zurich HB, 1 meter away."},
			{"Zurich", []Station{{"Zurich HB SZU", 2.7}, {"Zurich HB", 1.2}}, "The closest stations to Zurich are: Zurich HB SZU, 2 meters away; Zurich HB, 1 meter away."},
		},
		"de": []stationsTest{
			{"", []Station{}, "Ich konnte keine Haltestellen finden."},
			{"Zurich", []Station{}, "Ich konnte keine Haltestellen in der Nähe von Zurich finden."},
			{"", []Station{{"Zurich HB", 1.2}}, "Die nächste Haltestelle zu Ihnen ist: Zurich HB, 1 Meter entfernt."},
			{"Zurich", []Station{{"Zurich HB", 1.2}}, "Die nächste Haltestelle zum Zurich ist: Zurich HB, 1 Meter entfernt."},
			{"", []Station{{"Zurich HB SZU", 2.7}, {"Zurich HB", 1.2}}, "Die nächste Haltestellen zu Ihnen sind: Zurich HB SZU, 2 Meter entfernt; Zurich HB, 1 Meter entfernt."},
			{"Zurich", []Station{{"Zurich HB SZU", 2.7}, {"Zurich HB", 1.2}}, "Die nächste Haltestellen zum Zurich sind: Zurich HB SZU, 2 Meter entfernt; Zurich HB, 1 Meter entfernt."},
		},
	} {
		l := NewLocalizer(l, time.Now().Location())
		for _, want := range wants {
			got := l.Stations(want.Near, want.Stations)
			if got != want.Want {
				t.Errorf("want '%v', got '%v'", want.Want, got)
			}
		}
	}
}

type departuresTest struct {
	From      string
	To        string
	StartTime time.Time
	Deps      []Departure
	Want      string
}

func TestDepatures(t *testing.T) {
	for l, wants := range map[string][]departuresTest{
		"en": []departuresTest{
			{"Zurich", "", time.Time{}, []Departure{}, "I could not find any matching routes."},
			{"Zurich", "", time.Time{}, []Departure{{"S7", 0, "Zurich", "Enge", time.Unix(1517055015, 0), "bus", ""}},
				"The next departure from Zurich is: the S7 bus departing on-time at 12:10 to Enge."},
			{"Zurich", "", time.Unix(1517055015, 0), []Departure{{"S7", 0, "Zurich", "Enge", time.Unix(1517055015, 0), "bus", ""}},
				"The next departure leaving Zurich from 12:10 is: the S7 bus departing on-time at 12:10 to Enge."},
			{"Zurich", "", time.Time{}, []Departure{{"S7", 2, "Zurich", "Enge", time.Unix(1517055015, 0), "bus", ""}, {"S8", 0, "Zurich", "Basel", time.Unix(1517055015, 0), "train", "6"}},
				"The next 2 departures from Zurich are: the S7 bus departing at 12:10 with a 2-minute delay to Enge, and the S8 train departing on-time from platform 6 at 12:10 to Basel."},
			{"Zurich", "", time.Unix(1517055015, 0), []Departure{{"S7", 0, "Zurich", "Enge", time.Unix(1517055015, 0), "bus", ""}, {"S8", 0, "Zurich", "Basel", time.Unix(1517055015, 0), "train", ""}},
				"The next 2 departures leaving Zurich from 12:10 are: the S7 bus departing on-time at 12:10 to Enge, and the S8 train departing on-time at 12:10 to Basel."},
		},
		"de": []departuresTest{},
	} {
		l := NewLocalizer(l, time.Now().Location())
		for _, want := range wants {
			got := l.NextDepartures(want.From, want.To, want.StartTime, want.Deps)
			if got != want.Want {
				t.Errorf("want '%v', got '%v'", want.Want, got)
			}
		}
	}
}
