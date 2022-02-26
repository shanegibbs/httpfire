package agent

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ListenAddr: "0.0.0.0:8080",
	}
}

func Main(serverConfig ServerConfig) error {
	ctx := context.Background()
	ctx, triggerShutdown := context.WithCancel(ctx)

	config := AgentConfig{
		URL:         "http://127.0.0.1:8080",
		Timeout:     1 * time.Second,
		ThreadCount: 4,
		LogRequests: true,
	}

	{
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigterm
			triggerShutdown()
		}()
	}

	agent := NewAgent(ctx, triggerShutdown, config)

	return RunAgentServer(ctx, serverConfig, agent)
}

func RunAgentServer(ctx context.Context, serverConfig ServerConfig, agent *Agent) error {

	mux := http.NewServeMux()

	r := prometheus.NewRegistry()
	r.MustRegister(opsTotalMetrics)
	r.MustRegister(latencyMetrics)
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
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

			agent.SetConfig(config)
			agent.Start()
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
	mux.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			agent.Stop()
			agent.Start()
			w.WriteHeader(200)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
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
	log.Println("Listening...")

	<-ctx.Done()
	log.Print("Server shutting down...")

	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
