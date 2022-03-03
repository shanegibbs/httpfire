package main

import (
	"context"
	"log"
	"strconv"

	"github.com/shanegibbs/httpfire/pkg/agent"
	"github.com/shanegibbs/httpfire/pkg/common"
	"github.com/shanegibbs/httpfire/pkg/director"
)

func main() {
	ctx := context.Background()
	ctx, shutdown := context.WithCancel(ctx)
	common.HandleSigterm(shutdown)

	go func() {
		err := director.Main(ctx, shutdown, director.DefaultServerConfig())
		if err != nil {
			log.Printf("error: %v", err)
		}
		shutdown()
	}()

	startAgent := func(port uint64) {
		config := agent.ServerConfig{ListenAddr: "127.0.0.1:" + strconv.FormatUint(port, 10)}
		err := agent.Main(ctx, shutdown, config)
		if err != nil {
			log.Printf("error: %v", err)
		}
		shutdown()
	}

	go startAgent(8081)
	go startAgent(8082)

	<-ctx.Done()
}
