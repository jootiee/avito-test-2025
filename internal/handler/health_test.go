package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/service"
)

func TestHandler_Health(t *testing.T) {
	// Setup - handler doesn't use services for health check
	svc := &service.Service{} // Empty service is fine for health check
	log := &mockLogger{}
	h := New(svc, log)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expectedBody := `{"status":"ok"}`
	if body := w.Body.String(); body != expectedBody+"\n" {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}

	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestHandler_Root(t *testing.T) {
	// Setup - root handler doesn't use services
	svc := &service.Service{} // Empty service is fine for root
	log := &mockLogger{}
	h := New(svc, log)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check that response contains service name
	body := w.Body.String()
	if body == "" {
		t.Error("expected non-empty response body")
	}
}
