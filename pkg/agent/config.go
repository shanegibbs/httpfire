package agent

import "time"

type ServerConfig struct {
	ListenAddr string
}

type AgentConfig struct {
	Request            Request `json:"request"`
	ThreadCount        uint    `json:"threadCount"`
	LogRequests        bool    `json:"logRequests"`
	RateLimitPerSecond float64 `json:"rateLimitPerSecond"`
}

type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
}
