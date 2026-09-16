package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/matiasinsaurralde/sample-repo-go/pkg/config"
)

type helloRequest struct {
	Name string `json:"name"`
}

type helloResponse struct {
	Message string `json:"message"`
}

func helloHandler(c *gin.Context) {
	var req helloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if req.Name == "" {
		req.Name = "world"
	}

	c.JSON(http.StatusOK, helloResponse{
		Message: "hello " + req.Name,
	})
}

type lsResponse struct {
	Output string `json:"output"`
}

func lsHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	output, err := exec.Command("ls", path).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(string(output))})
		return
	}

	c.JSON(http.StatusOK, lsResponse{
		Output: strings.TrimRight(string(output), "\n"),
	})
}

type adminResponse struct {
	Addr string `json:"addr"`
}

// adminHandler reports non-sensitive service information. The token is
// supplied by configuration and compared in constant time; the response
// deliberately excludes the token and the process environment.
func adminHandler(cfg config.Config) gin.HandlerFunc {
	expected := sha256.Sum256([]byte(cfg.AdminToken))

	return func(c *gin.Context) {
		provided := sha256.Sum256([]byte(c.GetHeader("X-Admin-Token")))
		if subtle.ConstantTimeCompare(provided[:], expected[:]) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, adminResponse{Addr: cfg.Addr})
	}
}
