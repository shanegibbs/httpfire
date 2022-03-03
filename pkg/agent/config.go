package agent

import "time"

type ServerConfig struct {
	ListenAddr string
}

type AgentConfig struct {
	URL                string        `json:"url"`
	Timeout            time.Duration `json:"timeout"`
	ThreadCount        uint          `json:"threadCount"`
	LogRequests        bool          `json:"logRequests"`
	RateLimitPerSecond float64       `json:"rateLimitPerSecond"`
}
