package service

import (
	"testing"
)

func TestSelectRandomReviewers_EmptyCandidates(t *testing.T) {
	candidates := []string{}
	result := SelectRandomReviewers(candidates, 2)

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d reviewers", len(result))
	}
}

func TestSelectRandomReviewers_LessThanMax(t *testing.T) {
	candidates := []string{"user1", "user2"}
	result := SelectRandomReviewers(candidates, 5)

	if len(result) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(result))
	}

	// Verify all candidates are in result
	resultSet := make(map[string]bool)
	for _, r := range result {
		resultSet[r] = true
	}

	for _, c := range candidates {
		if !resultSet[c] {
			t.Errorf("expected %s in result", c)
		}
	}
}

func TestSelectRandomReviewers_ExactlyMax(t *testing.T) {
	candidates := []string{"user1", "user2", "user3"}
	result := SelectRandomReviewers(candidates, 3)

	if len(result) != 3 {
		t.Errorf("expected 3 reviewers, got %d", len(result))
	}
}

func TestSelectRandomReviewers_MoreThanMax(t *testing.T) {
	candidates := []string{"user1", "user2", "user3", "user4", "user5"}
	result := SelectRandomReviewers(candidates, 2)

	if len(result) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(result))
	}

	// Verify all returned reviewers are from candidates
	candidatesSet := make(map[string]bool)
	for _, c := range candidates {
		candidatesSet[c] = true
	}

	for _, r := range result {
		if !candidatesSet[r] {
			t.Errorf("returned reviewer %s not in candidates", r)
		}
	}

	// Verify no duplicates
	seen := make(map[string]bool)
	for _, r := range result {
		if seen[r] {
			t.Errorf("duplicate reviewer %s in result", r)
		}
		seen[r] = true
	}
}

func TestSelectRandomReviewers_MaxZero(t *testing.T) {
	candidates := []string{"user1", "user2", "user3"}
	result := SelectRandomReviewers(candidates, 0)

	if len(result) != 0 {
		t.Errorf("expected 0 reviewers for max=0, got %d", len(result))
	}
}

func TestFilterActiveCandidates_AllActive(t *testing.T) {
	users := []string{"user1", "user2", "user3"}
	activeMap := map[string]bool{
		"user1": true,
		"user2": true,
		"user3": true,
	}

	result := FilterActiveCandidates(users, activeMap)

	if len(result) != 3 {
		t.Errorf("expected 3 active users, got %d", len(result))
	}
}

func TestFilterActiveCandidates_SomeInactive(t *testing.T) {
	users := []string{"user1", "user2", "user3", "user4"}
	activeMap := map[string]bool{
		"user1": true,
		"user2": false,
		"user3": true,
		"user4": false,
	}

	result := FilterActiveCandidates(users, activeMap)

	if len(result) != 2 {
		t.Errorf("expected 2 active users, got %d", len(result))
	}

	// Verify only active users are returned
	for _, userID := range result {
		if !activeMap[userID] {
			t.Errorf("inactive user %s should not be in result", userID)
		}
	}
}

func TestFilterActiveCandidates_WithExclusions(t *testing.T) {
	users := []string{"user1", "user2", "user3", "user4"}
	activeMap := map[string]bool{
		"user1": true,
		"user2": true,
		"user3": true,
		"user4": true,
	}

	result := FilterActiveCandidates(users, activeMap, "user1", "user3")

	if len(result) != 2 {
		t.Errorf("expected 2 users after exclusions, got %d", len(result))
	}

	// Verify excluded users are not in result
	for _, userID := range result {
		if userID == "user1" || userID == "user3" {
			t.Errorf("excluded user %s should not be in result", userID)
		}
	}
}

func TestFilterActiveCandidates_ExcludeInactiveAndSpecific(t *testing.T) {
	users := []string{"user1", "user2", "user3", "user4", "user5"}
	activeMap := map[string]bool{
		"user1": true,
		"user2": false,
		"user3": true,
		"user4": true,
		"user5": false,
	}

	result := FilterActiveCandidates(users, activeMap, "user3")

	if len(result) != 2 {
		t.Errorf("expected 2 users (active, not excluded), got %d", len(result))
	}

	// Verify result contains only user1 and user4
	expectedSet := map[string]bool{"user1": true, "user4": true}
	for _, userID := range result {
		if !expectedSet[userID] {
			t.Errorf("unexpected user %s in result", userID)
		}
	}
}

func TestFilterActiveCandidates_EmptyResult(t *testing.T) {
	users := []string{"user1", "user2"}
	activeMap := map[string]bool{
		"user1": false,
		"user2": false,
	}

	result := FilterActiveCandidates(users, activeMap)

	if len(result) != 0 {
		t.Errorf("expected 0 users when all inactive, got %d", len(result))
	}
}

func TestFilterActiveCandidates_UserNotInActiveMap(t *testing.T) {
	users := []string{"user1", "user2", "user3"}
	activeMap := map[string]bool{
		"user1": true,
		"user2": true,
		// user3 not in map
	}

	result := FilterActiveCandidates(users, activeMap)

	if len(result) != 2 {
		t.Errorf("expected 2 users (only those in active map), got %d", len(result))
	}

	// user3 should not be in result as it's not in activeMap
	for _, userID := range result {
		if userID == "user3" {
			t.Error("user3 should not be in result when not in activeMap")
		}
	}
}

func TestFilterActiveCandidates_EmptyInput(t *testing.T) {
	users := []string{}
	activeMap := map[string]bool{
		"user1": true,
	}

	result := FilterActiveCandidates(users, activeMap)

	if len(result) != 0 {
		t.Errorf("expected 0 users for empty input, got %d", len(result))
	}
}
