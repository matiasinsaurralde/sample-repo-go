package api

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
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

	output, err := exec.Command("sh", "-c", "ls "+path).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": strings.TrimSpace(string(output))})
		return
	}

	c.JSON(http.StatusOK, lsResponse{
		Output: strings.TrimRight(string(output), "\n"),
	})
}
