package main

import (
	"log"

	"forum/config"
	"forum/internal/app"
)

func main() {
	cnf, err := config.New()
	if err != nil {
		log.Fatalf("Init Config Error: %v\n", err)
	}
	app.New(cnf).Start()
}
