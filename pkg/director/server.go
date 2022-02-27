package director

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shanegibbs/httpfire/pkg/agent"
)

func DefaultServerConfig() Config {
	return Config{
		ListenAddr: "0.0.0.0:8080",
	}
}

func Main(ctx context.Context, shutdown context.CancelFunc, config Config) error {
	log.Printf("using config: %v", config)

	director := NewDirector(ctx, shutdown, config)
	err := director.StartDiscovery()
	if err != nil {
		return fmt.Errorf("failed to start discovery: %v", err)
	}

	return RunDirectorServer(ctx, config, director)
}

func RunDirectorServer(ctx context.Context, config Config, director *Director) error {

	mux := http.NewServeMux()

	r := prometheus.NewRegistry()
	// r.MustRegister(opsTotalMetrics)
	// r.MustRegister(latencyMetrics)
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			config := agent.AgentConfig{}
			err := json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			director.Start(config)
			w.WriteHeader(200)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			director.Stop()
			w.WriteHeader(200)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	if config.ListenAddr == "" {
		return fmt.Errorf("listen_addr config is empty")
	}

	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: logRequestHandler(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Director listening on %v...", srv.Addr)

	<-ctx.Done()
	log.Print("Director shutting down...")

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
