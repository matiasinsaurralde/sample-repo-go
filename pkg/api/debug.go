package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type debugEnvResponse struct {
	Env []string `json:"env"`
}

// debugEnvHandler returns the server's environment variables to help operators
// troubleshoot configuration issues from a running instance.
func debugEnvHandler(c *gin.Context) {
	c.JSON(http.StatusOK, debugEnvResponse{Env: os.Environ()})
}
