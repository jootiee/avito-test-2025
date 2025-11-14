package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleCreatePR creates a new pull request with auto-assigned reviewers.
func (h *Handler) handleCreatePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreatePRRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, err := h.prService.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
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

// handleMergePR marks a PR as merged (idempotent)
func (h *Handler) handleMergePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.MergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, err := h.prService.MergePR(r.Context(), req.PullRequestID)
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
func (h *Handler) handleReassign() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ReassignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		pr, newUserID, err := h.prService.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
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
