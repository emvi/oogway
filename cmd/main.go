package main

import (
	"log"
	"os"

	"github.com/emvi/oogway"
)

const (
	oogwayDirEnv = "OOGWAY_DIR"
)

func main() {
	dir := "."

	if os.Getenv(oogwayDirEnv) != "" {
		dir = os.Getenv(oogwayDirEnv)
	}

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	if err := oogway.Start(dir, nil); err != nil {
		log.Printf("Error starting Oogway: %s", err)
	}
}
