package github

import "time"

var (
	// This is for the random-assignee rotation.
	reviewerRotation = map[string]struct{}{
		"slevenick":   {},
		"c2thorn":     {},
		"rileykarson": {},
		"melinath":    {},
		"ScottSuarez": {},
		"shuyama1":    {},
		"SarahFrench": {},
		"roaks3":      {},
		"zli82016":    {},
		"trodge":      {},
		"hao-nan-li":  {},
		"NickElliot":  {},
		"BBBmau":      {},
	}

	// This is for new team members who are onboarding
	trustedContributors = map[string]struct{}{}

	// This is for reviewers who are "on vacation": will not receive new review assignments but will still receive re-requests for assigned PRs.
	// User can specify the time zone like this, and following the example below:
	pdtLoc, _           = time.LoadLocation("America/Los_Angeles")
	bstLoc, _           = time.LoadLocation("Europe/London")
	onVacationReviewers = []onVacationReviewer{
		// Example: taking vacation from 2024-03-28 to 2024-04-02 in pdt time zone.
		// both ends are inclusive:
		// {
		// 	id:        "xyz",
		// 	startDate: newDate(2024, 3, 28, pdtLoc),
		// 	endDate:   newDate(2024, 4, 2, pdtLoc),
		// },
		{
			id:        "BBBmau",
			startDate: newDate(2024, 9, 26, pdtLoc),
			endDate:   newDate(2024, 10, 2, pdtLoc),
		},
		{
			id:        "hao-nan-li",
			startDate: newDate(2024, 9, 24, pdtLoc),
			endDate:   newDate(2024, 10, 4, pdtLoc),
		},
		{
			id:        "ScottSuarez",
			startDate: newDate(2024, 4, 30, pdtLoc),
			endDate:   newDate(2024, 7, 31, pdtLoc),
		},
		{
			id:        "shuyama1",
			startDate: newDate(2024, 9, 26, pdtLoc),
			endDate:   newDate(2024, 10, 4, pdtLoc),
		},
		{
			id:        "melinath",
			startDate: newDate(2024, 9, 18, pdtLoc),
			endDate:   newDate(2024, 9, 23, pdtLoc),
		},
		{
			id:        "slevenick",
			startDate: newDate(2024, 7, 5, pdtLoc),
			endDate:   newDate(2024, 7, 16, pdtLoc),
		},
		{
			id:        "c2thorn",
			startDate: newDate(2024, 7, 10, pdtLoc),
			endDate:   newDate(2024, 7, 16, pdtLoc),
		},
		{
			id:        "rileykarson",
			startDate: newDate(2024, 7, 18, pdtLoc),
			endDate:   newDate(2024, 8, 10, pdtLoc),
		},
		{
			id:        "roaks3",
			startDate: newDate(2024, 8, 2, pdtLoc),
			endDate:   newDate(2024, 8, 9, pdtLoc),
		},
		{
			id:        "slevenick",
			startDate: newDate(2024, 8, 10, pdtLoc),
			endDate:   newDate(2024, 8, 17, pdtLoc),
		},
		{
			id:        "trodge",
			startDate: newDate(2024, 8, 24, pdtLoc),
			endDate:   newDate(2024, 9, 2, pdtLoc),
		},
		{
			id:        "roaks3",
			startDate: newDate(2024, 9, 13, pdtLoc),
			endDate:   newDate(2024, 9, 20, pdtLoc),
		},
		{
			id:        "SarahFrench",
			startDate: newDate(2024, 9, 20, bstLoc),
			endDate:   newDate(2024, 9, 23, bstLoc),
		},
		{
			id:        "c2thorn",
			startDate: newDate(2024, 10, 2, bstLoc),
			endDate:   newDate(2024, 10, 14, bstLoc),
		},
	}
)
