package localize

import (
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/text/language"
)

const dataDir = "./data"

type Localizer struct {
	lang language.Tag
	tz   *time.Location
	t    i18n.TranslateFunc
}

var matcher = language.NewMatcher([]language.Tag{
	language.English, // The first language is used as fallback.
	language.German,
	language.French,
})

func init() {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".all.json") {
			i18n.MustLoadTranslationFile(path.Join(dataDir, f.Name()))
		}
	}
}

func NewLocalizer(lang string, timezone *time.Location) Localizer {
	tag, _ := language.MatchStrings(matcher, lang)
	// Use the tag to avoid https://github.com/nicksnyder/go-i18n/issues/76.
	b, _ := tag.Base()
	t := i18n.MustTfunc(b.String())
	return Localizer{tag, timezone, t}
}

type Station struct {
	Name     string
	Distance float64
}

type Departure struct {
	Name         string
	MinutesDelay int
	From         string
	To           string
	Departing    time.Time
	Mode         string
	Platform     string
}

func (l *Localizer) NeedLocation() string {
	return l.t("location_needed")
}

func (l *Localizer) PermissionContext() string {
	return l.t("to_look_for_stations")
}

func (l *Localizer) Stations(near string, stations []Station) string {
	parts := []string{}
	for _, s := range stations {
		parts = append(parts, l.t("meters_away", int(s.Distance), map[string]interface{}{"Name": s.Name}))
	}
	if len(parts) == 0 {
		if len(near) > 0 {
			return l.t("no_nearby_stations_near", map[string]interface{}{"Near": near})
		}
		return l.t("no_nearby_stations")
	}
	// XXX: This string join is bad i18n, but it works.
	if len(near) > 0 {
		return l.t("closest_near", len(parts), map[string]interface{}{"Near": near, "Stations": strings.Join(parts, "; ")})
	}
	return l.t("closest_to_you", len(parts), map[string]interface{}{"Stations": strings.Join(parts, "; ")})
}

func (l *Localizer) NextDepartures(from, to string, startTime time.Time, deps []Departure) string {
	parts := []string{}
	for _, d := range deps {
		// "the 7 tram departing on-time at 15:04 to Farbhof"
		// d.Name, d.Mode, d.MinutesDelay, d.Departing, d.MinutesDelay, d.To
		tm := d.Departing.In(l.tz).Format("15:04")
		var name string
		switch d.Mode {
		case "bus":
			name = l.t("bus", map[string]interface{}{"Name": d.Name})
		case "tram":
			name = l.t("tram", map[string]interface{}{"Name": d.Name})
		case "train":
			name = l.t("train", map[string]interface{}{"Name": d.Name})
		case "ship":
			name = l.t("ship", map[string]interface{}{"Name": d.Name})
		default:
			name = l.t("unknown_mode", map[string]interface{}{"Name": d.Name})
		}
		if d.Platform == "" {
			if d.MinutesDelay < 1 {
				parts = append(parts, l.t("the_7_tram_on_time_at_1504_to_farbhof", map[string]interface{}{
					"Name":        name,
					"Time":        tm,
					"Destination": d.To,
				}))
			} else {
				parts = append(parts, l.t("the_7_tram_with_a_5_minute_delay_at_1504_to_farbhof", map[string]interface{}{
					"Name":        name,
					"Time":        tm,
					"Destination": d.To,
					"Delay":       d.MinutesDelay,
				}))
			}
		} else {
			if d.MinutesDelay < 1 {
				parts = append(parts, l.t("the_7_tram_on_time_from_platform_2_at_1504_to_farbhof", map[string]interface{}{
					"Name":        name,
					"Time":        tm,
					"Destination": d.To,
					"Platform":    d.Platform,
				}))
			} else {
				parts = append(parts, l.t("the_7_tram_with_a_5_minute_delay_from_platform_2_at_1504_to_farbhof", map[string]interface{}{
					"Name":        name,
					"Time":        tm,
					"Destination": d.To,
					"Delay":       d.MinutesDelay,
					"Platform":    d.Platform,
				}))
			}
		}
	}

	if len(parts) == 0 {
		return l.t("could_not_find_any_routes")
	}

	if startTime.IsZero() && to == "" {
		return l.t("next_departures", len(parts), map[string]interface{}{
			"From":       from,
			"Departures": strings.Join(parts[:len(parts)-1], "; "),
			"Last":       parts[len(parts)-1],
		})
	} else if startTime.IsZero() && to != "" {
		return l.t("next_departures_to", len(parts), map[string]interface{}{
			"From":       from,
			"To":         to,
			"Departures": strings.Join(parts[:len(parts)-1], "; "),
			"Last":       parts[len(parts)-1],
		})
	} else if !startTime.IsZero() && to == "" {
		return l.t("next_departures_at", len(parts), map[string]interface{}{
			"From":       from,
			"Departures": strings.Join(parts[:len(parts)-1], "; "),
			"Last":       parts[len(parts)-1],
			"Time":       startTime.In(l.tz).Format("15:04"),
		})
	} else /* !startTime.IsZero() && to != "" */ {
		return l.t("next_departures_to_at", len(parts), map[string]interface{}{
			"From":       from,
			"To":         to,
			"Departures": strings.Join(parts[:len(parts)-1], "; "),
			"Last":       parts[len(parts)-1],
			"Time":       startTime.In(l.tz).Format("15:04"),
		})
	}
}
