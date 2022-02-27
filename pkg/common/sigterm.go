package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func HandleSigterm(shutdown context.CancelFunc) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigterm
		shutdown()
	}()
}
