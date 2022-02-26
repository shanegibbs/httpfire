package agent

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AgentConfig struct {
	URL         string
	Timeout     time.Duration
	ThreadCount uint
	LogRequests bool
}

func Main() error {
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

	return RunAgentServer(ctx, agent)
}

func RunAgentServer(ctx context.Context, agent *Agent) error {

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
		agent.Start()
		w.WriteHeader(200)
	})
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		agent.Stop()
		w.WriteHeader(200)
	})
	mux.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
		agent.Stop()
		agent.Start()
		w.WriteHeader(200)
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
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

		// call the original http.Handler we're wrapping
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method
		log.Println(method, uri)
	}

	// http.HandlerFunc wraps a function so that it
	// implements http.Handler interface
	return http.HandlerFunc(fn)
}
