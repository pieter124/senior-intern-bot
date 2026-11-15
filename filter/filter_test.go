package filter

import (
	domain "senior_intern_bot/domain"
	"testing"
)

func TestKeywordFilter_Classify(t *testing.T) {
	// An axis a case leaves unset is the zero Axis{}, which always evaluates
	// to unknown, so each case lists only the keywords that matter to it.
	var tests = []struct {
		name       string
		role       Axis
		region     Axis
		discipline Axis
		posting    domain.Posting
		want       domain.Verdict
	}{
		{
			name:       "all three axes in scope is accepted",
			role:       Axis{InScope: []string{"summer analyst"}},
			region:     Axis{InScope: []string{"berlin"}},
			discipline: Axis{InScope: []string{"technology"}},
			posting:    domain.Posting{Title: "Summer Analyst 2027 Technology", Location: "Berlin"},
			want:       domain.ACCEPT,
		},
		{
			name:       "unknown discipline does not block accept",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"london"}},
			discipline: Axis{InScope: []string{"software"}, OutOfScope: []string{"marketing"}},
			posting:    domain.Posting{Title: "Internship", Location: "London"},
			want:       domain.ACCEPT,
		},
		{
			name:       "out-of-scope role is rejected",
			role:       Axis{InScope: []string{"internship"}, OutOfScope: []string{"senior"}},
			region:     Axis{InScope: []string{"london"}},
			discipline: Axis{InScope: []string{"software"}},
			posting:    domain.Posting{Title: "Senior Software Engineer", Location: "London"},
			want:       domain.REJECT,
		},
		{
			name:       "out-of-scope region is rejected",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"london"}, OutOfScope: []string{"san francisco"}},
			discipline: Axis{InScope: []string{"software"}},
			posting:    domain.Posting{Title: "Software Internship", Location: "San Francisco"},
			want:       domain.REJECT,
		},
		{
			name:       "out-of-scope discipline is rejected",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"london"}},
			discipline: Axis{InScope: []string{"software"}, OutOfScope: []string{"marketing"}},
			posting:    domain.Posting{Title: "Marketing Internship", Location: "London"},
			want:       domain.REJECT,
		},
		{
			name:       "unknown region goes to review",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"london"}},
			discipline: Axis{InScope: []string{"software"}},
			posting:    domain.Posting{Title: "Software Internship", Location: "Remote"},
			want:       domain.REVIEW,
		},
		{
			name:       "unknown role goes to review",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"europe"}},
			discipline: Axis{InScope: []string{"software"}},
			posting:    domain.Posting{Title: "Software Engineer", Location: "Europe"},
			want:       domain.REVIEW,
		},
		{
			name:       "region in and out of scope contradicts itself and goes to review",
			role:       Axis{InScope: []string{"internship"}},
			region:     Axis{InScope: []string{"london"}, OutOfScope: []string{"san francisco"}},
			discipline: Axis{InScope: []string{"software"}},
			posting:    domain.Posting{Title: "Software Internship", Location: "London or San Francisco"},
			want:       domain.REVIEW,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := New(tt.role, tt.region, tt.discipline)
			got := filter.Classify(tt.posting)
			if got != tt.want {
				t.Errorf("Classify(%+v): got %+v, want %+v", tt.posting, got, tt.want)
			}
		})
	}
}
