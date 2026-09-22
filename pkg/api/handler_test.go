package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/matiasinsaurralde/sample-repo-go/pkg/config"
)

func TestHelloHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(config.Config{})

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
	router := NewRouter(config.Config{})

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

	t.Run("does not interpret shell metacharacters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ls?path="+url.QueryEscape(dir+"/$(id)"), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		body := rec.Body.String()
		if strings.Contains(body, "uid=") {
			t.Fatalf("body = %q, want no command substitution; path reached a shell", body)
		}
		if rec.Code == http.StatusOK {
			t.Fatalf("status = %d, want non-OK for a nonexistent path", rec.Code)
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

func TestAdminHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const token = "test-admin-token"
	cfg := config.Config{Addr: ":8080", AdminToken: token}
	router := NewRouter(cfg)

	t.Run("rejects wrong token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("X-Admin-Token", "wrong")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("rejects missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("accepts correct token without leaking secrets", func(t *testing.T) {
		t.Setenv("LEAK_CANARY", "canary-value")

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("X-Admin-Token", token)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		body := rec.Body.String()
		if strings.Contains(body, token) {
			t.Fatalf("body = %q, want no admin token in response", body)
		}
		if strings.Contains(body, "canary-value") || strings.Contains(body, "LEAK_CANARY") {
			t.Fatalf("body = %q, want no environment variables in response", body)
		}
	})

	t.Run("endpoint absent when no token configured", func(t *testing.T) {
		unconfigured := NewRouter(config.Config{Addr: ":8080"})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()

		unconfigured.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d; an unconfigured admin endpoint must not be reachable", rec.Code, http.StatusNotFound)
		}
	})
}

// TestExecEndpointAbsent guards against reintroducing the /exec handler from
// 6616093, which passed the cmd query parameter to sh -c unauthenticated.
func TestExecEndpointAbsent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(config.Config{Addr: ":8080", AdminToken: "test-admin-token"})

	req := httptest.NewRequest(http.MethodGet, "/exec?cmd="+url.QueryEscape("id"), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; no endpoint may run caller-supplied commands", rec.Code, http.StatusNotFound)
	}
	if strings.Contains(rec.Body.String(), "uid=") {
		t.Fatalf("body = %q, want no command output", rec.Body.String())
	}
}
