package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestLsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "alpha.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "beta.txt"), []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("lists directory", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ls?path="+dir, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "alpha.txt") || !strings.Contains(body, "beta.txt") {
			t.Fatalf("body = %q, want listing of alpha.txt and beta.txt", body)
		}
	})

	t.Run("missing path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ls", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if rec.Body.String() != `{"error":"path query parameter is required"}` {
			t.Fatalf("body = %q, want path required error", rec.Body.String())
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ls?path=/no/such/directory", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}
