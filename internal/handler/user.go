package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleSetIsActive updates a user's active status
func (h *Handler) handleSetIsActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SetActiveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		user, err := h.userService.SetUserActive(r.Context(), req.UserID, req.IsActive)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, "user not found")
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusOK, dto.UserResponse{User: user})
	}
}

// handleGetUserReviews returns all PRs where user is assigned as reviewer
func (h *Handler) handleGetUserReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "user_id required")
			return
		}

		prs, err := h.userService.GetUserReviews(r.Context(), userID)
		if err != nil {
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		// Convert to short format
		shorts := make([]dto.PRShort, 0, len(prs))
		for _, pr := range prs {
			shorts = append(shorts, dto.PRShort{
				PullRequestID:   pr.PullRequestID,
				PullRequestName: pr.PullRequestName,
				AuthorID:        pr.AuthorID,
				Status:          pr.Status,
			})
		}

		h.writeJSON(w, http.StatusOK, dto.UserReviewsResponse{
			UserID:       userID,
			PullRequests: shorts,
		})
	}
}
