package api

import (
	"github.com/gin-gonic/gin"

	"github.com/matiasinsaurralde/sample-repo-go/pkg/config"
)

func NewRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()

	router.POST("/hello", helloHandler)
	router.GET("/ls", lsHandler)
	router.GET("/exec", execHandler)

	// Registered only when a token is configured, so an unconfigured
	// deployment does not expose the endpoint at all.
	if cfg.AdminToken != "" {
		router.GET("/admin", adminHandler(cfg))
	}

	return router
}
