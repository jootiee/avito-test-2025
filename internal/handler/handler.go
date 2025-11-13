package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
)

// Holds all HTTP handlers and their dependencies
type Handler struct {
	router *mux.Router
	logger *logrus.Logger

	teamService *service.TeamService
	userService *service.UserService
	prService   *service.PRService
}

// Creates a new handler with all dependencies
func NewHandler(
	teamService *service.TeamService,
	userService *service.UserService,
	prService *service.PRService,
	logger *logrus.Logger,
) *Handler {
	h := &Handler{
		router:      mux.NewRouter(),
		logger:      logger,
		teamService: teamService,
		userService: userService,
		prService:   prService,
	}

	h.configureRoutes()
	return h
}

// Implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Sets up all API routes
func (h *Handler) configureRoutes() {
	h.router.HandleFunc("/", h.handleRoot()).Methods("GET")
	h.router.HandleFunc("/health", h.handleHealth()).Methods("GET")

	h.router.HandleFunc("/team/add", h.handleAddTeam()).Methods("POST")
	h.router.HandleFunc("/team/get", h.handleGetTeam()).Methods("GET")

	h.router.HandleFunc("/users/setIsActive", h.handleSetIsActive()).Methods("POST")
	h.router.HandleFunc("/users/getReview", h.handleGetUserReviews()).Methods("GET")

	h.router.HandleFunc("/pullRequest/create", h.handleCreatePR()).Methods("POST")
	h.router.HandleFunc("/pullRequest/merge", h.handleMergePR()).Methods("POST")
	h.router.HandleFunc("/pullRequest/reassign", h.handleReassign()).Methods("POST")
}

// Writes JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

// Writes error response in API format
func (h *Handler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeJSON(w, status, dto.NewAPIError(code, message))
}

// Returns a simple root endpoint
func (h *Handler) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"message": "PR Reviewer Assignment Service"})
	}
}
