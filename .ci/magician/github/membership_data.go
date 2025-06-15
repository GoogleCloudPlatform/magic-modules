package github

var (
	// This is for the random-assignee rotation.
	reviewerRotation = ReviewerRotation{
		"BBBmau": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 4, 7),
					endDate:   newDate(2025, 4, 11),
				},
			},
		},
		"c2thorn": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 4, 9),
					endDate:   newDate(2025, 4, 15),
				},
			},
		},
		"hao-nan-li": {
			vacations: []Vacation{},
		},
		"melinath": {
			vacations: []Vacation{},
		},
		"NickElliot": {
			vacations: []Vacation{},
		},
		"rileykarson": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 2, 25),
					endDate:   newDate(2025, 3, 10),
				},
			},
		},
		"roaks3": {
			vacations: []Vacation{},
		},
		"ScottSuarez": {
			vacations: []Vacation{},
		},
		"shuyama1": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 5, 23),
					endDate:   newDate(2025, 5, 30),
				},
			},
		},
		"SirGitsalot": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 1, 18),
					endDate:   newDate(2025, 1, 25),
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
			vacations: []Vacation{},
		},
		"zli82016": {
			vacations: []Vacation{
				{
					startDate: newDate(2025, 1, 15),
					endDate:   newDate(2025, 2, 9),
				},
			},
		},
	}

	// This is for new team members who are onboarding
	trustedContributors = map[string]struct{}{
		"bbasata":           struct{}{},
		"jaylonmcshan03":    struct{}{},
		"malhotrasagar2212": struct{}{},
	}
)
