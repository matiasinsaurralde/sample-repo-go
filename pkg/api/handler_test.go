package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHelloHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "with name",
			body:       `{"name":"world"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello world"}`,
		},
		{
			name:       "default name when omitted",
			body:       `{}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello world"}`,
		},
		{
			name:       "custom name",
			body:       `{"name":"alice"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello alice"}`,
		},
		{
			name:       "invalid json",
			body:       `{not json}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid JSON body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
