package github

import (
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	start1 := time.Now().Add(-8 * 24 * time.Hour)
	end1 := time.Now().Add(-4 * 24 * time.Hour)
	rr := ReviewerRotation{
		"id1": {
			vacations: []Vacation{
				{
					start: start1,
					end:   end1,
				},
			},
		},
		"id2": {
			vacations: []Vacation{
				{
					start: start1,
					end:   end1,
				},
			},
		},
	}
	data := []byte(`id1:
  vacations:
  - start: YYYY/MM/DD
    end: YYYY/MM/DD
id2:
  timezone: US/Eastern
  vacations:
  - start: 1969/07/16
    end: 1969/07/24
`)
	if err := rr.read(data); err != nil {
		t.Fatal(err)
	}
	if len(rr) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(rr))
	}
	if len(rr["id1"].vacations) != 1 {
		// Confirm sample vacations are not loaded.
		t.Fatalf("expected 1 vacation for id1, got %d", len(rr["id1"].vacations))
	}
	if !rr["id1"].vacations[0].start.Equal(start1) || !rr["id1"].vacations[0].end.Equal(end1) {
		t.Fatalf("expected id1's vacation start %s and end1 %s, got %s and %s", start1, end1, rr["id1"].vacations[0].start, rr["id1"].vacations[0].end)
	}
	if len(rr["id2"].vacations) != 2 {
		// Confirm vacations merge.
		t.Fatalf("expected 2 vacations for id2, got %d", len(rr["id2"].vacations))
	}
	ny, _ := time.LoadLocation("America/New_York")
	start2, _ := time.ParseInLocation("2006/01/02", "1969/07/16", ny)
	if !rr["id2"].vacations[1].start.Equal(start2) {
		t.Fatalf("expected id2's vacation start %s, got %s", start2, rr["id2"].vacations[1].start)
	}
}
