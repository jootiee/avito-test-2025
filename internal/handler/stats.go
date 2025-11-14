package handler

import (
	"net/http"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleGetStats godoc
// @Summary Get reviewer assignment statistics
// @Description Returns statistics about reviewer assignments per user and per pull request
// @Tags Statistics
// @Produce json
// @Success 200 {object} dto.StatsResponse "Statistics"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /stats [get]
func (h *Handler) handleGetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userCounts, prCounts, err := h.service.PullRequest.GetStatistics(ctx)
		if err != nil {
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, "Failed to retrieve statistics")
			return
		}

		userStats := make([]dto.UserStats, 0, len(userCounts))
		for userID, count := range userCounts {
			userStats = append(userStats, dto.UserStats{
				UserID:          userID,
				AssignmentCount: count,
			})
		}

		prStats := make([]dto.PullRequestStats, 0, len(prCounts))
		for prID, count := range prCounts {
			prStats = append(prStats, dto.PullRequestStats{
				PullRequestID: prID,
				ReviewerCount: count,
			})
		}

		response := dto.StatsResponse{
			UserStats: userStats,
			PRStats:   prStats,
		}

		h.writeJSON(w, http.StatusOK, response)
	}
}
