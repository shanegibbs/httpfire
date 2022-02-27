package main

import (
	"context"
	"log"

	"github.com/shanegibbs/httpfire/pkg/common"
	"github.com/shanegibbs/httpfire/pkg/director"
)

func main() {
	ctx := context.Background()
	ctx, shutdown := context.WithCancel(ctx)
	common.HandleSigterm(shutdown)

	config, err := director.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	err = director.Main(ctx, shutdown, *config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
