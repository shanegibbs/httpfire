package main

import (
	"log"

	"github.com/shanegibbs/httpfire/pkg/director"
)

func main() {
	err := director.Main(director.DefaultServerConfig())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
