package director

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/shanegibbs/httpfire/pkg/agent"
)

type Director struct {
	ctx      context.Context
	shutdown context.CancelFunc
	config   ServerConfig
	client   *http.Client
	agents   []agent.Agent
}

func NewDirector(ctx context.Context, shutdown context.CancelFunc, config ServerConfig) *Director {
	client := http.DefaultClient
	return &Director{ctx, shutdown, config, client, []agent.Agent{}}
}

func (d *Director) AddAgent(endpoint string) error {
	agent, err := agent.NewRemoteAgent(d.ctx, d.client, endpoint)
	if err != nil {
		return fmt.Errorf("failed to create agent (%s): %v", endpoint, err)
	}
	d.agents = append(d.agents, agent)
	return nil
}

func (d *Director) Start(config agent.AgentConfig) {
	for _, agent := range d.agents {
		err := agent.Start(config)
		if err != nil {
			log.Printf("failed to start agent: %v", err)
		}
	}
}

func (d *Director) Stop() {
	for _, agent := range d.agents {
		err := agent.Stop()
		if err != nil {
			log.Printf("failed to stop agent: %v", err)
		}
	}
}
