package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Agent interface {
	Start(config AgentConfig) error
	Stop() error
}

type LocalAgent struct {
	ctx         context.Context
	agentAbort  context.CancelFunc
	stopSession context.CancelFunc
	config      AgentConfig
}

func NewAgent(ctx context.Context, cancel context.CancelFunc) *LocalAgent {
	return &LocalAgent{
		ctx:         ctx,
		agentAbort:  cancel,
		stopSession: nil,
	}
}

func (a *LocalAgent) Start(config AgentConfig) error {
	if a.stopSession != nil {
		log.Println("agent already running")
		return nil
	}

	a.config = config

	workerCtx, stopSession := context.WithCancel(a.ctx)
	a.stopSession = stopSession

	go func() {
		log.Println("Starting")
		for i := 0; i < int(a.config.ThreadCount); i++ {
			id := uint(i)
			go func() {
				err := a.WorkerFunc(id, workerCtx)

				if errors.Is(err, context.Canceled) {
					log.Printf("worker %v canceled", id)

				} else {
					log.Printf("workerfunc err: %v", err)
					// abort program when there is a worker error
					a.agentAbort()
				}
			}()
		}
	}()
	return nil
}

func (a *LocalAgent) Stop() error {
	log.Println("Stopping")
	if a.stopSession != nil {
		a.stopSession()
		a.stopSession = nil
	}
	return nil
}

var (
	opsTotalMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "httpfire_ops_total",
	}, []string{"status"})
	latencyMetrics = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "httpfire_latency",
		Buckets: prometheus.ExponentialBuckets(20, 1.3, 25),
	})
)

func (a *LocalAgent) calculateWaitTime(requestDuration time.Duration) time.Duration {
	if a.config.RateLimitPerSecond == 0 {
		return time.Duration(0)
	}

	nanoSecondsPerRequst := (1.0 / a.config.RateLimitPerSecond) * 1000000000.0
	durationPerRequest := time.Duration(nanoSecondsPerRequst)
	return durationPerRequest - requestDuration
}

func (a *LocalAgent) WorkerFunc(id uint, ctx context.Context) error {
	var i uint
	client := &http.Client{}

	log.Printf("Starting worker func loop %v", id)

	for {
		start := time.Now()

		if err := a.ExecuteOperation(ctx, id, i, client); err != nil {
			return err
		}

		i++

		select {
		case <-time.After(a.calculateWaitTime(time.Since(start))):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *LocalAgent) ExecuteOperation(ctx context.Context, id, i uint, client *http.Client) error {
	reqConfig := a.config.Request

	reqCtx, cancel := context.WithTimeout(ctx, reqConfig.Timeout)
	defer cancel()

	body := strings.NewReader(reqConfig.Body)
	req, err := http.NewRequestWithContext(reqCtx, reqConfig.Method, reqConfig.URL, body)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "httpfire/0.0")
	for k, v := range reqConfig.Headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	res, err := client.Do(req)
	latency := time.Since(start)

	msgPrefix := fmt.Sprintf("worker(%v).%v:", id, i)
	statusCode := 0
	msg := ""

	if errors.Is(err, context.Canceled) {
		return err
	} else if err != nil {
		msg = fmt.Sprintf("error: %v", err)
	} else {
		msg = fmt.Sprintf("%s %s %v", reqConfig.URL, res.Status, latency)
		statusCode = res.StatusCode
	}

	opsTotalMetrics.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	latencyMetrics.Observe(float64(latency.Milliseconds()))

	if a.config.LogRequests {
		log.Printf("%s %s", msgPrefix, msg)
	}

	return nil
}
