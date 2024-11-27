package main

import (
	"log"

	"github.com/demtoni/tade/internal/api"
	"github.com/demtoni/tade/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := api.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Run())
}
