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
	config   Config
	client   *http.Client
	agents   []agent.Agent
}

func NewDirector(ctx context.Context, shutdown context.CancelFunc, config Config) *Director {
	client := http.DefaultClient
	return &Director{ctx, shutdown, config, client, []agent.Agent{}}
}

func (d *Director) StartDiscovery() error {
	if d.config.Discovery == nil {
		return nil
	}

	discovery := *d.config.Discovery
	if discovery.Static != nil {
		d.setupStaticDiscovery(*discovery.Static)
	}

	return nil
}

func (d *Director) setupStaticDiscovery(discovery StaticDiscovery) error {
	for _, ep := range discovery.Endpoints {
		err := d.AddAgent(ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Director) AddAgent(endpoint string) error {
	log.Printf("Adding agent: %s", endpoint)
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
