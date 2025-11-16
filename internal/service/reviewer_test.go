package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReviewerTestSuite struct {
	suite.Suite
}

func (s *ReviewerTestSuite) TestSelectRandom() {
	type args struct {
		candidates []string
		max        int
	}

	tests := []struct {
		name       string
		args       args
		assertFunc func(t *testing.T, got []string, candidates []string, maxCount int)
	}{
		{
			name: "empty candidates",
			args: args{candidates: []string{}, max: 2},
			assertFunc: func(t *testing.T, got []string, _ []string, _ int) {
				assert.Len(t, got, 0)
			},
		},
		{
			name: "less than max returns all (permutation)",
			args: args{candidates: []string{"user1", "user2"}, max: 5},
			assertFunc: func(t *testing.T, got []string, candidates []string, _ int) {
				assert.Len(t, got, 2)
				assert.ElementsMatch(t, candidates, got)
			},
		},
		{
			name: "exactly max returns all (permutation)",
			args: args{candidates: []string{"user1", "user2", "user3"}, max: 3},
			assertFunc: func(t *testing.T, got []string, candidates []string, _ int) {
				assert.Len(t, got, 3)
				assert.ElementsMatch(t, candidates, got)
			},
		},
		{
			name: "more than max returns unique subset",
			args: args{candidates: []string{"user1", "user2", "user3", "user4", "user5"}, max: 2},
			assertFunc: func(t *testing.T, got []string, candidates []string, maxCount int) {
				assert.Len(t, got, maxCount)
				// ensure subset of candidates and no duplicates
				seen := map[string]struct{}{}
				valid := map[string]struct{}{}
				for _, c := range candidates {
					valid[c] = struct{}{}
				}
				for _, r := range got {
					if _, ok := valid[r]; !ok {
						t.Fatalf("returned reviewer %s not in candidates", r)
					}
					if _, dup := seen[r]; dup {
						t.Fatalf("duplicate reviewer %s in result", r)
					}
					seen[r] = struct{}{}
				}
			},
		},
		{
			name: "max zero returns empty",
			args: args{candidates: []string{"user1", "user2", "user3"}, max: 0},
			assertFunc: func(t *testing.T, got []string, _ []string, _ int) {
				assert.Empty(t, got)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			got := SelectRandomReviewers(tt.args.candidates, tt.args.max)
			tt.assertFunc(s.T(), got, tt.args.candidates, tt.args.max)
		})
	}
}

func (s *ReviewerTestSuite) TestFilterActiveCandidates() {
	type args struct {
		users     []string
		activeMap map[string]bool
		exclude   []string
	}

	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name: "all active",
			args: args{
				users:     []string{"user1", "user2", "user3"},
				activeMap: map[string]bool{"user1": true, "user2": true, "user3": true},
			},
			expected: []string{"user1", "user2", "user3"},
		},
		{
			name: "some inactive",
			args: args{
				users:     []string{"user1", "user2", "user3", "user4"},
				activeMap: map[string]bool{"user1": true, "user2": false, "user3": true, "user4": false},
			},
			expected: []string{"user1", "user3"},
		},
		{
			name: "with exclusions",
			args: args{
				users:     []string{"user1", "user2", "user3", "user4"},
				activeMap: map[string]bool{"user1": true, "user2": true, "user3": true, "user4": true},
				exclude:   []string{"user1", "user3"},
			},
			expected: []string{"user2", "user4"},
		},
		{
			name: "exclude inactive and specific",
			args: args{
				users:     []string{"user1", "user2", "user3", "user4", "user5"},
				activeMap: map[string]bool{"user1": true, "user2": false, "user3": true, "user4": true, "user5": false},
				exclude:   []string{"user3"},
			},
			expected: []string{"user1", "user4"},
		},
		{
			name: "empty result when all inactive",
			args: args{
				users:     []string{"user1", "user2"},
				activeMap: map[string]bool{"user1": false, "user2": false},
			},
			expected: []string{},
		},
		{
			name: "user not in active map excluded",
			args: args{
				users:     []string{"user1", "user2", "user3"},
				activeMap: map[string]bool{"user1": true, "user2": true},
			},
			expected: []string{"user1", "user2"},
		},
		{
			name: "empty input",
			args: args{
				users:     []string{},
				activeMap: map[string]bool{"user1": true},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			got := FilterActiveCandidates(tt.args.users, tt.args.activeMap, tt.args.exclude...)
			assert.ElementsMatch(s.T(), tt.expected, got)
		})
	}
}

func TestReviewer(t *testing.T) {
	suite.Run(t, new(ReviewerTestSuite))
}
