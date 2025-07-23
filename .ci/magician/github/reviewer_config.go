package github

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type ReviewerConfig struct {
	// timezone controls the timezone for vacation start / end dates. Default: US/Pacific.
	timezone *time.Location

	// vacations allows specifying times when new reviews should not be requested of the reviewer.
	// Existing PRs will still have reviews re-requested.
	// Both startDate and endDate are inclusive.
	vacations []Vacation
}

type reviewerConfigYAML struct {
	Timezone  string
	Vacations []vacationYAML
}

const defaultTimezone = "US/Pacific"

func (rc ReviewerConfig) MarshalYAML() (any, error) {
	yrc := reviewerConfigYAML{
		Timezone:  rc.timezone.String(),
		Vacations: make([]vacationYAML, len(rc.vacations)),
	}
	if rc.timezone == nil {
		yrc.Timezone = defaultTimezone
	}
	if len(rc.vacations) == 0 {
		yrc.Vacations = []vacationYAML{
			vacationYAML{
				Start: sampleDateStr,
				End:   sampleDateStr,
			},
		}
	}
	for i, v := range rc.vacations {
		marshalled, err := v.MarshalYAML()
		if err != nil {
			return nil, err
		}
		var ok bool
		yrc.Vacations[i], ok = marshalled.(vacationYAML)
		if !ok {
			return nil, fmt.Errorf("non-vacationYAML returned from MarshalYAML: %v (%T)", marshalled, marshalled)
		}
	}
	return yrc, nil
}

func (rc *ReviewerConfig) UnmarshalYAML(value *yaml.Node) error {
	var yrc reviewerConfigYAML
	if err := value.Decode(&yrc); err != nil {
		return fmt.Errorf("failed to decode reviewer config: %w", err)
	}

	timezoneStr := yrc.Timezone
	if timezoneStr == "" {
		timezoneStr = defaultTimezone
	}
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return fmt.Errorf("failed to load timezone %q: %w", timezoneStr, err)
	}
	rc.timezone = loc

	for _, vYAML := range yrc.Vacations {
		if vYAML.Start == sampleDateStr && vYAML.End == sampleDateStr {
			continue // Skip sample placeholder vacations
		}

		parsedStartDate, err := time.Parse("2006/01/02", vYAML.Start)
		if err != nil {
			return fmt.Errorf("failed to parse start date %q: %w", vYAML.Start, err)
		}
		parsedEndDate, err := time.Parse("2006/01/02", vYAML.End)
		if err != nil {
			return fmt.Errorf("failed to parse end date %q: %w", vYAML.End, err)
		}

		vacation := Vacation{
			// Date components are from parsed dates, time set to day boundaries in reviewer's timezone.
			start: time.Date(parsedStartDate.Year(), parsedStartDate.Month(), parsedStartDate.Day(), 0, 0, 0, 0, rc.timezone),
			end:   time.Date(parsedEndDate.Year(), parsedEndDate.Month(), parsedEndDate.Day(), 23, 59, 59, 0, rc.timezone),
		}

		rc.vacations = append(rc.vacations, vacation)
	}

	return nil
}

var usPacific, _ = time.LoadLocation("US/Pacific")

// Set start and end for each vacation based on startDate and endDate.
func (rc *ReviewerConfig) setStartEnd() {
	for i := range rc.vacations {
		v := &rc.vacations[i]
		loc := usPacific
		if rc.timezone != nil {
			loc = rc.timezone
		}
		if v.start.IsZero() && v.end.IsZero() {
			// Adjust vacation times to the reviewer's timezone
			v.start = time.Date(v.startDate.year, time.Month(v.startDate.month), v.startDate.day, 0, 0, 0, 0, loc)
			v.end = time.Date(v.endDate.year, time.Month(v.endDate.month), v.endDate.day, 23, 59, 59, 0, loc)
		}
	}
}

const (
	vacationStartOffset = -8 * time.Hour // Vacations will effectively start at 4pm the previous day in the reviewer's timezone instead of midnight.
	vacationEndOffset   = 9 * time.Hour  // Vacations will effectively end at 9am the next day in the reviewer's timezone instead of midnight.
)

func (rc *ReviewerConfig) onVacation(now time.Time) bool {
	for _, v := range rc.vacations {
		if v.start.Add(vacationStartOffset).Before(now) && v.end.Add(vacationEndOffset).After(now) {
			return true
		}
	}
	return false
}
