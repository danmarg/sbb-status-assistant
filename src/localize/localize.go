package localize

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type Localizer struct {
	lang language.Tag
	tz   *time.Location
}

var matcher = language.NewMatcher([]language.Tag{
	language.English, // The first language is used as fallback.
	language.German,
})

func NewLocalizer(lang string) Localizer {
	tag, _ := language.MatchStrings(matcher, lang)
	tz, _ := time.LoadLocation("Europe/Zurich")
	return Localizer{tag, tz}
}

type Station struct {
	Name     string
	Distance int
}

type Departure struct {
	Name      string
	OnTime    bool
	To        string
	Departing time.Time
	Mode      string
	Platform  string
}

func (l *Localizer) NeedLocation() string {
	if l.lang == language.German {
		return "Ich brauche Ihren Standort."
	}
	return "I need your location."
}

func (l *Localizer) Stations(near string, stations []Station) string {
	parts := []string{}
	for _, s := range stations {
		part := s.Name
		if s.Distance > 0 {
			if l.lang == language.German {
				part += fmt.Sprintf(", %d Meter entfernt", s.Distance)
			} else {
				part += fmt.Sprintf(", %d meters away", s.Distance)
			}
		}
	}
	if len(parts) == 0 {
		if l.lang == language.German {
			return fmt.Sprintf("Ich konnte keine Haltestellen in der Nähe von %s finden.", near)
		}
		fmt.Sprintf("I could not find any matching stations near %s.", near)
	} else if len(parts) == 1 {
		if l.lang == language.German {
			return fmt.Sprintf("Die nächste Haltestelle zum %s ist: %s,", near, parts[0])
		}
		return fmt.Sprintf("The closest station to %s is: %s,", parts[0])
	}
	if l.lang == language.German {
		return fmt.Sprintf("Die nächste Haltestellen zum %s sind: %s,", near, strings.Join(parts, ";"))
	}
	return fmt.Sprintf("The closest stations to %s are: %s,", near, strings.Join(parts, ";"))
}

func (l *Localizer) NextDepartures(from string, startTime time.Time, deps []Departure) string {
	parts := []string{}
	for _, d := range deps {
		// "the 7 tram departing on-time at 15:04 to Farbhof"
		var part string
		var mode string
		if l.lang == language.German {
			switch d.Mode {
			case "bus":
				part += "der "
				mode = "Bus"
			case "tram":
				part += "die "
				mode = "Tram"
			case "train":
				part += "der "
				mode = "Zug"
			case "ship":
				part += "das "
				mode = "Shiff"
			default:
				part += "das "
				mode = d.Mode
			}
		} else {
			part += "the "
		}
		part += fmt.Sprintf("%s %s ", d.Name, mode)
		tm := d.Departing.In(l.tz).Format("15:04")
		if l.lang == language.German {
			if d.OnTime {
				part += "pünktlich abfahren "
			} else {
				part += "abfahren mit einer Verspätung "
			}
			if d.Platform != "" {
				part += fmt.Sprintf("von Gleis %s ", d.Platform)
			}
			part += fmt.Sprintf("nach %s um %s", d.To, tm)
		} else {
			if d.OnTime {
				part += "departing on-time "
			} else {
				part += "departing behind schedule "
			}
			if d.Platform != "" {
				part += fmt.Sprintf("from platform %s ", d.Platform)
			}

			part += fmt.Sprintf("to %s at %s", d.To, tm)
		}
		parts = append(parts, part)
	}

	if len(parts) == 0 {
		if l.lang == language.German {
			return "Ich konnte keine passenden Haltestellen oder Linien finden."
		}
		return "I could not find any matching stations or routes."
	} else if len(parts) == 1 {
		if l.lang == language.German {
			return fmt.Sprintf("Die nächste Abfahrt von %s ist %s.", from, parts[0])
		} else {
			return fmt.Sprintf("The next departure from %s is %s.", from, parts[0])
		}
	}
	var result string
	if startTime.IsZero() {
		if l.lang == language.German {
			result = fmt.Sprintf("Die nächsten %d Abfahrten von %s sind: ", len(parts), from)
		} else {
			result = fmt.Sprintf("The next %d departures from %s are: ", len(parts), from)
		}
	} else {
		st := startTime.In(l.tz).Format("15:04")
		if l.lang == language.German {
			result = fmt.Sprintf("Die nächsten %d Abfahrten ab %s von %s sind: ", len(parts), st, from)
		} else {
			result = fmt.Sprintf("The next %d departures leaving %s from %s are: ", len(parts), from, st)
		}

	}

	result += strings.Join(parts[:len(parts)-1], "; ")
	if l.lang == language.German {
		result += " und "
	} else {
		result += " and "
	}
	return result + parts[len(parts)-1] + "."
}
