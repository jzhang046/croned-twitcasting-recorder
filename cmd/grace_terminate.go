package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const terminationGraceDuration = 3 * time.Second

func newInterruptableCtx() context.Context {
	rootCtx, cancelFunc := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt

		log.Printf("Terminating in %s.. \n", terminationGraceDuration)
		cancelFunc()
		time.Sleep(terminationGraceDuration)
	}()

	return rootCtx
}
