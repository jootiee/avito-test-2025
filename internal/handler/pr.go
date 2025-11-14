package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleCreatePR creates a new pull request with auto-assigned reviewers.
// @Summary Create a new pull request
// @Description Creates a new PR and automatically assigns up to 2 reviewers from the author's team
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param pr body dto.CreatePRRequest true "Pull request details"
// @Success 201 {object} dto.PRResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Author or team not found"
// @Failure 409 {object} dto.ErrorResponse "PR already exists"
// @Failure 500 {object} dto.ErrorResponse
// @Router /prs/create [post]
func (h *Handler) handleCreatePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePRRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, err := h.service.PR.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				h.writeError(w, http.StatusConflict, dto.ErrCodePRExists, "PR id already exists")
				return
			}
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusCreated, dto.PRResponse{PR: pr})
	}
}

// handleMergePR marks a PR as merged
// @Summary Merge a pull request
// @Description Marks a pull request as merged
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.MergeRequest true "PR ID to merge"
// @Success 200 {object} dto.PRResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "PR not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /prs/merge [post]
func (h *Handler) handleMergePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.MergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, err := h.service.PR.MergePR(r.Context(), req.PullRequestID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, "pr not found")
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusOK, dto.PRResponse{PR: pr})
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
// @Failure 404 {object} dto.ErrorResponse "PR or user not found"
// @Failure 409 {object} dto.ErrorResponse "PR merged, not assigned, or no candidate"
// @Failure 500 {object} dto.ErrorResponse
// @Router /prs/reassign [post]
func (h *Handler) handleReassign() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ReassignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, newUserID, err := h.service.PR.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
		if err != nil {
			if strings.Contains(err.Error(), "merged PR") {
				h.writeError(w, http.StatusConflict, dto.ErrCodePRMerged, "cannot reassign on merged PR")
				return
			}
			if strings.Contains(err.Error(), "not assigned") {
				h.writeError(w, http.StatusConflict, dto.ErrCodeNotAssigned, "reviewer is not assigned to this PR")
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
			PR:         pr,
			ReplacedBy: newUserID,
		})
	}
}
