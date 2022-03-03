package agent

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ListenAddr: "0.0.0.0:8080",
	}
}

func Main(ctx context.Context, shutdown context.CancelFunc, serverConfig ServerConfig) error {
	agent := NewAgent(ctx, shutdown)
	return RunAgentServer(ctx, serverConfig, agent)
}

func RunAgentServer(ctx context.Context, serverConfig ServerConfig, agent *LocalAgent) error {

	mux := http.NewServeMux()

	r := prometheus.NewRegistry()
	r.MustRegister(opsTotalMetrics)
	r.MustRegister(latencyMetrics)
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	})
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			config := AgentConfig{}
			err := json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			agent.Start(config)
			w.WriteHeader(200)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			agent.Stop()
			w.WriteHeader(200)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// mux.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "POST":
	// 		agent.Stop()
	// 		agent.Start()
	// 		w.WriteHeader(200)
	// 	default:
	// 		w.WriteHeader(http.StatusMethodNotAllowed)
	// 	}
	// })
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(agent.config)
			if err != nil {
				log.Printf("failed to write response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		Addr:    serverConfig.ListenAddr,
		Handler: logRequestHandler(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Agent listening on %v...", srv.Addr)

	<-ctx.Done()
	log.Print("Agent shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	return srv.Shutdown(shutdownCtx)
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		uri := r.URL.String()
		method := r.Method
		log.Println(method, uri)
	}
	return http.HandlerFunc(fn)
}
