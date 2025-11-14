package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/jootiee/avito-test-2025/docs" // Import generated docs
	"github.com/jootiee/avito-test-2025/internal/dto"
	"github.com/jootiee/avito-test-2025/internal/service"
	"github.com/jootiee/avito-test-2025/pkg/logger"
)

type Handler struct {
	router  *mux.Router
	logger  logger.Interface
	service *service.Service
}

func New(
	svc *service.Service,
	logger logger.Interface,
) *Handler {
	h := &Handler{
		router:  mux.NewRouter(),
		logger:  logger,
		service: svc,
	}

	h.configureRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) configureRoutes() {
	h.router.Use(h.RequestIDMiddleware)
	h.router.Use(h.LoggingMiddleware)

	h.router.HandleFunc("/", h.handleRoot()).Methods("GET")
	h.router.HandleFunc("/health", h.handleHealth()).Methods("GET")

	// Swagger documentation
	h.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	h.router.HandleFunc("/team/add", h.handleAddTeam()).Methods("POST")
	h.router.HandleFunc("/team/get", h.handleGetTeam()).Methods("GET")

	h.router.HandleFunc("/users/setIsActive", h.handleSetIsActive()).Methods("POST")
	h.router.HandleFunc("/users/getReview", h.handleGetUserReviews()).Methods("GET")

	h.router.HandleFunc("/pullRequest/create", h.handleCreatePR()).Methods("POST")
	h.router.HandleFunc("/pullRequest/merge", h.handleMergePR()).Methods("POST")
	h.router.HandleFunc("/pullRequest/reassign", h.handleReassign()).Methods("POST")

	h.router.HandleFunc("/stats", h.handleGetStats()).Methods("GET")
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeJSON(w, status, dto.NewAPIError(code, message))
}

func (h *Handler) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"message": "PR Reviewer Assignment Service"})
	}
}
