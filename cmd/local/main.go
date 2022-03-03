package main

import (
	"context"
	"time"

	"github.com/shanegibbs/httpfire/pkg/agent"
	"github.com/shanegibbs/httpfire/pkg/common"
)

func main() {
	ctx := context.Background()
	ctx, shutdown := context.WithCancel(ctx)
	common.HandleSigterm(shutdown)

	config := agent.AgentConfig{
		Request: agent.Request{
			URL:    "http://localhost:9090",
			Host:   "example.com",
			Method: "GET",
			Headers: map[string]string{
				"HttpFire-Session": "1",
			},
			Body:    "",
			Timeout: 1 * time.Second,
		},
		ThreadCount:        1,
		RateLimitPerSecond: 0.5,
		LogRequests:        true,
	}

	agent := agent.NewAgent(ctx, shutdown)
	agent.Start(config)

	<-ctx.Done()
}
