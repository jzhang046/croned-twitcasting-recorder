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

func newInterruptableCtx() (context.Context, <-chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	rootCtx, cancelOnInterrupt := context.WithCancel(context.Background())
	afterGraceTermination := make(chan struct{})

	go func() {
		<-interrupt

		log.Printf("Terminating in %s.. \n", terminationGraceDuration)
		cancelOnInterrupt()
		time.Sleep(terminationGraceDuration)

		close(afterGraceTermination)
	}()

	return rootCtx, afterGraceTermination
}
