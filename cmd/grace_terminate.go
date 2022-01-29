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

func waitForInterruput(cancelFunc context.CancelFunc) <-chan struct{} {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	interrupted := make(chan struct{})

	go func() {
		<-interrupt

		log.Printf("Terminating in %s.. \n", terminationGraceDuration)
		go cancelFunc()

		time.Sleep(terminationGraceDuration)
		close(interrupted)
	}()

	return interrupted
}
