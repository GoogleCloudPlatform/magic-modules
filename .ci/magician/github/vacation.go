package github

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type Vacation struct {
	start, end         time.Time
	startDate, endDate date
}

type vacationYAML struct {
	Start, End string
}

const sampleDateStr = "YYYY/MM/DD"

type date struct {
	year, month, day int
}

func newDate(year, month, day int) date {
	return date{year: year, month: month, day: day}
}

func (v Vacation) MarshalYAML() (any, error) {
	return vacationYAML{
		Start: v.start.Format("2006/01/02"),
		End:   v.end.Format("2006/01/02"),
	}, nil
}

func (v *Vacation) UnmarshalYAML(value *yaml.Node) error {
	var vacationMap map[string]string
	if err := value.Decode(&vacationMap); err != nil {
		return err
	}

	startDateStr, ok := vacationMap["start"]
	if !ok {
		return fmt.Errorf("vacation missing start date")
	}
	endDateStr, ok := vacationMap["end"]
	if !ok {
		return fmt.Errorf("vacation missing end date")
	}

	if startDateStr == sampleDateStr && endDateStr == sampleDateStr {
		return nil
	}

	startDate, err := time.Parse("2006/01/02", startDateStr)
	if err != nil {
		return fmt.Errorf("failed to parse start date %q: %w", startDateStr, err)
	}
	endDate, err := time.Parse("2006/01/02", endDateStr)
	if err != nil {
		return fmt.Errorf("failed to parse end date %q: %w", endDateStr, err)
	}

	v.start = startDate
	v.end = endDate
	return nil
}
