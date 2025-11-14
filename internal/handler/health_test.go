package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jootiee/avito-test-2025/internal/service"
)

func TestHandler_Health(t *testing.T) {
	svc := &service.Service{}
	log := &mockLogger{}
	h := New(svc, log)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expectedBody := `{"status":"ok"}`
	if body := w.Body.String(); body != expectedBody+"\n" {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestHandler_Root(t *testing.T) {
	svc := &service.Service{}
	log := &mockLogger{}
	h := New(svc, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected non-empty response body")
	}
}
