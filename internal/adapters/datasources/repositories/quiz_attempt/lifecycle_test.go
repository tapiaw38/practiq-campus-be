package quiz_attempt

import (
	"testing"
	"time"
)

func secs(n int) *int { return &n }

func TestDeadlineFor(t *testing.T) {
	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC)
	closesIn10 := now.Add(10 * time.Minute)
	closesIn90 := now.Add(90 * time.Minute)

	cases := []struct {
		name   string
		params StartParams
		want   *time.Time
	}{
		{
			name:   "untimed quiz with no closing time has no deadline",
			params: StartParams{},
		},
		{
			name:   "time limit alone counts from the start",
			params: StartParams{TimeLimitSecs: secs(3600)},
			want:   ptr(now.Add(time.Hour)),
		},
		{
			// A whole hour's allowance cannot outlast the quiz itself: a
			// student starting ten minutes before it closes gets ten minutes,
			// not an hour reaching past the close.
			name:   "the quiz's closing time caps a longer allowance",
			params: StartParams{TimeLimitSecs: secs(3600), HardDeadline: &closesIn10},
			want:   &closesIn10,
		},
		{
			name:   "an allowance ending sooner than the close wins",
			params: StartParams{TimeLimitSecs: secs(3600), HardDeadline: &closesIn90},
			want:   ptr(now.Add(time.Hour)),
		},
		{
			// Untimed, but the quiz still closes: submitting after it closed
			// was possible because only the start checked availability.
			name:   "closing time applies even when there is no time limit",
			params: StartParams{HardDeadline: &closesIn10},
			want:   &closesIn10,
		},
		{
			name:   "a zero or negative limit means untimed, not instantly over",
			params: StartParams{TimeLimitSecs: secs(0)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := deadlineFor(now, tc.params)
			switch {
			case tc.want == nil && got != nil:
				t.Fatalf("deadlineFor = %v, want no deadline", *got)
			case tc.want != nil && got == nil:
				t.Fatalf("deadlineFor = no deadline, want %v", *tc.want)
			case tc.want != nil && !got.Equal(*tc.want):
				t.Fatalf("deadlineFor = %v, want %v", *got, *tc.want)
			}
		})
	}
}

func ptr(t time.Time) *time.Time { return &t }
