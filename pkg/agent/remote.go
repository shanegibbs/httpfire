package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
)

type RemoteAgent struct {
	ctx      context.Context
	client   *http.Client
	endpoint *url.URL
}

func NewRemoteAgent(ctx context.Context, client *http.Client, endpoint string) (Agent, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent endpoint URL %s: %v", endpoint, err)
	}

	return &RemoteAgent{ctx, client, url}, nil
}

func (a *RemoteAgent) Start(config AgentConfig) error {
	log.Printf("Sending start command to %v", a.endpoint)

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	startURL := *a.endpoint
	startURL.Path = path.Join(startURL.Path, "start")

	req, err := http.NewRequestWithContext(a.ctx, "POST", startURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("agent responded with failure (%s): %v", startURL.String(), res.Status)
	}

	return nil
}

func (a *RemoteAgent) Stop() error {
	log.Printf("Sending stop command to %v", a.endpoint)
	return nil
}
