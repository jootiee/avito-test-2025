package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jootiee/avito-test-2025/internal/domain"
	"github.com/jootiee/avito-test-2025/internal/dto"
)

// handleAddTeam creates a new team with members
// @Summary Add a new team
// @Description Creates a new team with specified members
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body dto.TeamAddRequest true "Team details"
// @Success 201 {object} dto.TeamResponse
// @Failure 400 {object} dto.ErrorResponse "Team already exists"
// @Failure 500 {object} dto.ErrorResponse
// @Router /teams/add [post]
func (h *Handler) handleAddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.TeamAddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "invalid request body")
			return
		}

		members := make([]domain.User, 0, len(req.Members))
		for _, m := range req.Members {
			members = append(members, domain.User{
				UserID:   m.UserID,
				Username: m.Username,
				TeamName: req.TeamName,
				IsActive: m.IsActive,
			})
		}

		team, err := h.service.Team.CreateTeam(r.Context(), req.TeamName, members)
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
// @Summary Get team details
// @Description Retrieves team information including all members
// @Tags Teams
// @Produce json
// @Param team_name query string true "Team name"
// @Success 200 {object} dto.TeamResponse
// @Failure 400 {object} dto.ErrorResponse "Missing team_name"
// @Failure 404 {object} dto.ErrorResponse "Team not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /teams/get [get]
func (h *Handler) handleGetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			h.writeError(w, http.StatusBadRequest, dto.ErrCodeNotFound, "team_name required")
			return
		}

		team, err := h.service.Team.GetTeam(r.Context(), teamName)
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
