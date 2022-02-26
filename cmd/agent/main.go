package main

import (
	"log"

	"github.com/shanegibbs/httpfire/pkg/agent"
)

func main() {
	err := agent.Main()
	if err != nil {
		log.Fatalf("error: %v", agent.Main())
	}
}
