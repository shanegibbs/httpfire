package agent

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type RemoteAgent struct {
	ctx      context.Context
	client   *http.Client
	endpoint string
}

func NewRemoteAgent(ctx context.Context, client *http.Client, endpoint string) Agent {
	return &RemoteAgent{ctx, client, endpoint}
}

func (a *RemoteAgent) Start(config AgentConfig) error {
	log.Printf("Sending start command to %v", a.endpoint)

	req, err := http.NewRequestWithContext(a.ctx, "POST", a.endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("agent responded with failure: %v", err)
	}

	return nil
}

func (a *RemoteAgent) Stop() error {
	log.Printf("Sending stop command to %v", a.endpoint)
	return nil
}
