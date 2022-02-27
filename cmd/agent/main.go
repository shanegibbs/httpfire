package main

import (
	"context"
	"log"

	"github.com/shanegibbs/httpfire/pkg/agent"
	"github.com/shanegibbs/httpfire/pkg/common"
)

func main() {
	ctx := context.Background()
	ctx, shutdown := context.WithCancel(ctx)
	common.HandleSigterm(shutdown)

	err := agent.Main(ctx, shutdown, agent.DefaultServerConfig())
	if err != nil {
		log.Printf("error: %v", err)
	}
}
