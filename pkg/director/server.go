package director

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shanegibbs/httpfire/pkg/agent"
)

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ListenAddr: "0.0.0.0:8080",
		AgentEndpoints: []string{
			"http://127.0.0.1:8081",
			"http://127.0.0.1:8082",
		},
	}
}

func Main(ctx context.Context, shutdown context.CancelFunc, serverConfig ServerConfig) error {
	director := NewDirector(ctx, shutdown, serverConfig)

	for _, ep := range serverConfig.AgentEndpoints {
		err := director.AddAgent(ep)
		if err != nil {
			return err
		}
	}

	return RunDirectorServer(ctx, serverConfig, director)
}

func RunDirectorServer(ctx context.Context, serverConfig ServerConfig, director *Director) error {

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

	srv := &http.Server{
		Addr:    serverConfig.ListenAddr,
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
