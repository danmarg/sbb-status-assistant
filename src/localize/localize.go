package localize

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type Localizer struct {
	Lang language.Tag
}

var matcher = language.NewMatcher([]language.Tag{
	language.English, // The first language is used as fallback.
	language.German,
})

func NewLocalizer(lang string) Localizer {
	tag, _ := language.MatchStrings(matcher, lang)
	return Localizer{tag}
}

type Departure struct {
	Name      string
	OnTime    bool
	To        string
	Departing time.Time
	Mode      string
}

func (l *Localizer) NextDepartures(from string, startTime time.Time, deps []Departure) string {
	parts := []string{}
	for _, d := range deps {
		// "the 7 tram departing on-time at 15:04 to Farbhof"
		var part string
		var mode string
		if l.Lang == language.German {
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
		part += fmt.Sprintf("%s %s", d.Name, mode)
		tm := d.Departing.Format("15:04")
		if l.Lang == language.German {
			if d.OnTime {
				part += fmt.Sprintf(" pünktlich abfahren nach %s um %s", d.To, tm)
			} else {
				part += fmt.Sprintf(" abfahren nach %s mit einer Verspätung um %s", d.To, tm)
			}
		} else {
			if d.OnTime {
				part += fmt.Sprintf(" departing on-time to %s at %s", d.To, tm)
			} else {
				part += fmt.Sprintf(" departing behind schedule to %s at %s", d.To, tm)
			}
		}
		parts = append(parts, part)
	}

	if len(parts) == 0 {
		if l.Lang == language.German {
			return "Ich konnte keine passenden Haltestellen oder Linien."
		}
		return "I could not find any matching stations or routes."
	} else if len(parts) == 1 {
		if l.Lang == language.German {
			return fmt.Sprintf("Die nächste Abfahrt von %s ist %s.", from, parts[0])
		} else {
			return fmt.Sprintf("The next departure from %s is %s.", from, parts[0])
		}
	}
	var result string
	if startTime.IsZero() {
		if l.Lang == language.German {
			result = fmt.Sprintf("Die nächsten %d Abfahrten von %s sind: ", len(parts), from)
		} else {
			result = fmt.Sprintf("The next %d departures from %s are: ", len(parts), from)
		}
	} else {
		st := startTime.Format("15:04")
		if l.Lang == language.German {
			result = fmt.Sprintf("Die nächsten %d Abfahrten ab %s von %s sind: ", len(parts), st, from)
		} else {
			result = fmt.Sprintf("The next %d departures leaving %s from %s are: ", len(parts), from, st)
		}

	}

	result += strings.Join(parts[:len(parts)-1], "; ")
	if l.Lang == language.German {
		result += " und "
	} else {
		result += " and "
	}
	return result + parts[len(parts)-1] + "."
}
