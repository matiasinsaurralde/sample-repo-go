package main

import (
	"log"

	"github.com/matiasinsaurralde/sample-repo-go/pkg/api"
	"github.com/matiasinsaurralde/sample-repo-go/pkg/config"
)

func main() {
	cfg := config.Load()

	router := api.NewRouter()
	log.Printf("listening on %s", cfg.Addr)
	if err := router.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
