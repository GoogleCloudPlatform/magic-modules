package github

import "time"

type date struct {
	year  int
	month int
	day   int
}

func newDate(year, month, day int) date {
	return date{
		year:  year,
		month: month,
		day:   day,
	}
}

type Vacation struct {
	startDate, endDate date
}

// GetStart returns a time corresponding to the beginning of the start date in the given timezone.
func (v Vacation) GetStart(timezone *time.Location) time.Time {
	if timezone == nil {
		timezone = usPacific
	}
	return time.Date(v.startDate.year, time.Month(v.startDate.month), v.startDate.day, 0, 0, 0, 0, timezone)
}

// GetEnd returns a time corresponding to the end of the end date in the given timezone
func (v Vacation) GetEnd(timezone *time.Location) time.Time {
	if timezone == nil {
		timezone = usPacific
	}
	return time.Date(v.endDate.year, time.Month(v.endDate.month), v.endDate.day, 0, 0, 0, 0, timezone).AddDate(0, 0, 1).Add(-1 * time.Millisecond)
}

type ReviewerConfig struct {
	// timezone controls the timezone for vacation start / end dates. Default: US/Pacific.
	timezone *time.Location

	// vacations allows specifying times when new reviews should not be requested of the reviewer.
	// Existing PRs will still have reviews re-requested.
	// Both startDate and endDate are inclusive.
	// Example: taking vacation from 2024-03-28 to 2024-04-02.
	// {
	// 	 vacations:        []Vacation{
	//     startDate: newDate(2024, 3, 28),
	// 	   endDate:   newDate(2024, 4, 2),
	//   },
	// },
	vacations []Vacation
}

var (
	usPacific, _ = time.LoadLocation("US/Pacific")
	usCentral, _ = time.LoadLocation("US/Central")
	usEastern, _ = time.LoadLocation("US/Eastern")
	london, _    = time.LoadLocation("Europe/London")

	// This is for the random-assignee rotation.
	reviewerRotation = map[string]ReviewerConfig{
		"BBBmau": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 7, 1),
					endDate:   newDate(2025, 7, 17),
				},
			},
		},
		"c2thorn": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 12, 3),
					endDate:   newDate(2025, 12, 15),
				},
			},
		},
		"hao-nan-li": {
			vacations: []Vacation{},
		},
		"malhotrasagar2212": {
			vacations: []Vacation{},
		},
		"melinath": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 12, 18),
					endDate:   newDate(2026, 1, 6),
				},
			},
		},
		"NickElliot": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 10, 22),
					endDate:   newDate(2025, 10, 28),
				},
			},
		},
		"rileykarson": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 2, 25),
					endDate:   newDate(2025, 3, 10),
				},
				{
					startDate: newDate(2025, 11, 11),
					endDate:   newDate(2025, 11, 24),
				},
			},
		},
		"roaks3": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 8, 1),
					endDate:   newDate(2025, 8, 11),
				},
				{
					startDate: newDate(2025, 9, 23),
					endDate:   newDate(2025, 9, 29),
				},
			},
		},
		"ScottSuarez": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 12, 5),
					endDate:   newDate(2025, 12, 15),
				},
			},
		},
		"shuyama1": {
			vacations: []Vacation{
				{
					startDate: newDate(2026, 01, 29),
					endDate:   newDate(2026, 02, 04),
				},
			},
		},
		"SirGitsalot": {
			vacations: []Vacation{
				{
					startDate: newDate(2026, 1, 18),
					endDate:   newDate(2026, 1, 27),
				},
			},
		},
		"slevenick": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 5, 22),
					endDate:   newDate(2025, 6, 7),
				},
			},
		},
		"trodge": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 8, 7),
					endDate:   newDate(2025, 8, 10),
				},
				{
					startDate: newDate(2025, 9, 18),
					endDate:   newDate(2025, 9, 28),
				},
			},
		},
		"zli82016": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 8, 27),
					endDate:   newDate(2025, 9, 2),
				},
			},
		},
	}

	// This is for new team members who are onboarding
	trustedContributors = map[string]struct{}{
		"bbasata": struct{}{},
	}
)
