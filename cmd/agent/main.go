package main

import (
	"log"

	"github.com/shanegibbs/httpfire/pkg/agent"
)

func main() {
	err := agent.Main(agent.DefaultServerConfig())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
