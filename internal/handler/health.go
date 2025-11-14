package handler

import "net/http"

// handleHealth returns service health status
func (h *Handler) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
