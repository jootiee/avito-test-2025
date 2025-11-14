package handler

import "net/http"

// handleHealth returns service health status
// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
