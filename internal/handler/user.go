package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleSetIsActive updates a user's active status
// @Summary Set user active status
// @Description Updates whether a user is active/available for PR review assignment
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.SetActiveRequest true "User active status"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/setIsActive [post]
func (h *Handler) handleSetIsActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SetActiveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		user, err := h.service.User.SetActive(r.Context(), req.UserID, req.IsActive)
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
// @Summary Get user's assigned PRs
// @Description Returns all pull requests where the user is assigned as a reviewer
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} dto.UserReviewsResponse
// @Failure 400 {object} dto.ErrorResponse "Missing user_id"
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/getReview [get]
func (h *Handler) handleGetUserReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "user_id required")
			return
		}

		prs, err := h.service.User.GetReviews(r.Context(), userID)
		if err != nil {
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

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
