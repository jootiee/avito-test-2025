package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleAddTeam creates a new team with members
func (h *Handler) handleAddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.TeamAddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		// Convert DTO to domain models
		members := make([]domain.User, 0, len(req.Members))
		for _, m := range req.Members {
			members = append(members, domain.User{
				UserID:   m.UserID,
				Username: m.Username,
				TeamName: req.TeamName,
				IsActive: m.IsActive,
			})
		}

		// Call service
		team, err := h.teamService.CreateTeam(r.Context(), req.TeamName, members)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				h.writeError(w, http.StatusBadRequest, dto.ErrCodeTeamExists, "team_name already exists")
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusCreated, dto.TeamResponse{Team: team})
	}
}

// handleGetTeam retrieves a team by name
func (h *Handler) handleGetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "team_name required")
			return
		}

		team, err := h.teamService.GetTeam(r.Context(), teamName)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				h.writeError(w, http.StatusNotFound, dto.ErrCodeNotFound, "team not found")
				return
			}
			h.writeError(w, http.StatusInternalServerError, dto.ErrCodeNotFound, err.Error())
			return
		}

		h.writeJSON(w, http.StatusOK, team)
	}
}
