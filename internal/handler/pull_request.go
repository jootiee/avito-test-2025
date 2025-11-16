package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleCreatePullRequest creates a new pull request with auto-assigned reviewers.
// @Summary Create a new pull request
// @Description Creates a new Pull Request and automatically assigns up to 2 reviewers from the author's team
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param pullRequest body dto.CreatePullRequestRequest true "Pull request details"
// @Success 201 {object} dto.PullRequestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Author or team not found"
// @Failure 409 {object} dto.ErrorResponse "PullRequest already exists"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/create [post]
func (h *Handler) handleCreatePullRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePullRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pullRequest, err := h.service.PullRequest.Create(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				h.writeError(w, http.StatusConflict, dto.ErrCodePullRequestExists, "PullRequest id already exists")
				return
			}
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusCreated, dto.PullRequestResponse{PullRequest: pullRequest})
	}
}

// handleMergePullRequest marks a PullRequest as merged
// @Summary Merge a pull request
// @Description Marks a pull request as merged
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.MergeRequest true "PullRequest ID to merge"
// @Success 200 {object} dto.PullRequestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "PullRequest not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/merge [post]
func (h *Handler) handleMergePullRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.MergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pullRequest, err := h.service.PullRequest.Merge(r.Context(), req.PullRequestID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, "PullRequest not found")
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusOK, dto.PullRequestResponse{PullRequest: pullRequest})
	}
}

// handleReassign replaces one reviewer with another from the same team
// @Summary Reassign a reviewer
// @Description Replaces a reviewer with another active member from the same team
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.ReassignRequest true "Reassignment details"
// @Success 200 {object} dto.ReassignResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "PullRequest or user not found"
// @Failure 409 {object} dto.ErrorResponse "PullRequest merged, not assigned, or no candidate"
// @Failure 500 {object} dto.ErrorResponse
// @Router /pullRequest/reassign [post]
func (h *Handler) handleReassign() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ReassignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pullRequest, newUserID, err := h.service.PullRequest.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
		if err != nil {
			if strings.Contains(err.Error(), "merged PullRequest") {
				h.writeError(w, http.StatusConflict, dto.ErrCodePullRequestMerged, "cannot reassign on merged PullRequest")
				return
			}
			if strings.Contains(err.Error(), "not assigned") {
				h.writeError(w, http.StatusConflict, dto.ErrCodeNotAssigned, "reviewer is not assigned to this PullRequest")
				return
			}
			if strings.Contains(err.Error(), "no active replacement") {
				h.writeError(w, http.StatusConflict, dto.ErrCodeNoCandidate, "no active replacement candidate in team")
				return
			}
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusOK, dto.ReassignResponse{
			PullRequest: pullRequest,
			ReplacedBy:  newUserID,
		})
	}
}
