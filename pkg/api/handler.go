package api

import (
	"net/http"

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
