package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jootiee/avito-test-2025/internal/service"
)

type HealthHandlerTestSuite struct {
	suite.Suite
}

func (s *HealthHandlerTestSuite) TestHealthEndpoints() {
	svc := &service.Service{}
	h := New(svc, &mockLogger{})

	tests := []struct {
		name       string
		path       string
		assertFunc func(w *httptest.ResponseRecorder)
	}{
		{
			name: "health ok",
			path: "/health",
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				assert.Equal(s.T(), "application/json", w.Header().Get("Content-Type"))
				assert.Equal(s.T(), "{\"status\":\"ok\"}\n", w.Body.String())
			},
		},
		{
			name: "root page",
			path: "/",
			assertFunc: func(w *httptest.ResponseRecorder) {
				assert.Equal(s.T(), http.StatusOK, w.Code)
				assert.NotEmpty(s.T(), w.Body.String())
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, tc.path, http.NoBody)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			tc.assertFunc(w)
		})
	}
}

func TestHealthHandler(t *testing.T) {
	suite.Run(t, new(HealthHandlerTestSuite))
}
