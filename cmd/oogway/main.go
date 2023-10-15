package main

import (
	oogway "github.com/emvi/oogway/pkg"
	"log"
	"os"
	"strings"
)

const (
	oogwayDirEnv = "OOGWAY_DIR"
)

func main() {
	cmd := "run"
	dir := "."

	if os.Getenv(oogwayDirEnv) != "" {
		dir = os.Getenv(oogwayDirEnv)
	}

	if len(os.Args) > 1 {
		cmd = strings.ToLower(os.Args[1])
	}

	if len(os.Args) > 2 {
		dir = os.Args[2]
	}

	switch cmd {
	case "run":
		if err := oogway.Start(dir, nil); err != nil {
			log.Printf("Error starting Oogway: %s", err)
		}
	case "init":
		if err := oogway.Init(dir); err != nil {
			log.Printf("Error initializing new project: %s", err)
		}
	default:
		log.Printf("Command unknown. Usage: oogway run|init <path>")
	}
}
