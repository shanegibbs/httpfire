package director

import (
	"context"
	"log"
	"net/http"

	"github.com/shanegibbs/httpfire/pkg/agent"
)

type Director struct {
	ctx      context.Context
	shutdown context.CancelFunc
	config   ServerConfig
	client   *http.Client
}

func NewDirector(ctx context.Context, shutdown context.CancelFunc, config ServerConfig) *Director {
	return &Director{ctx, shutdown, config, http.DefaultClient}
}

func (d *Director) Start(config agent.AgentConfig) {
	for _, ep := range d.config.AgentEndpoints {
		agent := agent.NewRemoteAgent(d.ctx, d.client, ep)
		err := agent.Start(config)
		if err != nil {
			log.Printf("failed to start agent (%s): %v", ep, err)
		}
	}
}

func (d *Director) Stop() {
	for _, ep := range d.config.AgentEndpoints {
		agent := agent.NewRemoteAgent(d.ctx, d.client, ep)
		agent.Stop()
	}
}
