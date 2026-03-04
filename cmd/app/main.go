package main

import (
	"log"

	"github.com/KimNattanan/go-chat-backend/internal/app"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig("")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
