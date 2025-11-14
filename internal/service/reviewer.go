package service

import (
	"math/rand"
)

// SelectRandomReviewers selects up to max reviewers randomly from candidates
func SelectRandomReviewers(candidates []string, max int) []string {
	if len(candidates) == 0 {
		return []string{}
	}
	if len(candidates) <= max {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		return candidates
	}
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	return candidates[:max]
}

// FilterActiveCandidates returns only active users, excluding specified IDs
func FilterActiveCandidates(users []string, activeMap map[string]bool, excludeIDs ...string) []string {
	excluded := make(map[string]struct{})
	for _, id := range excludeIDs {
		excluded[id] = struct{}{}
	}

	var candidates []string
	for _, userID := range users {
		if _, skip := excluded[userID]; skip {
			continue
		}
		if active, ok := activeMap[userID]; ok && active {
			candidates = append(candidates, userID)
		}
	}
	return candidates
}
