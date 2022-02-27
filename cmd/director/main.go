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

	err := director.Main(ctx, shutdown, director.DefaultServerConfig())
	if err != nil {
		log.Printf("error: %v", err)
	}
}
