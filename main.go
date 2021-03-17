package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
)

const (
	testStreamer = "kaguramea_vov"
)

func main() {
	log.Println("croned recorder started ")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every 3m", recordFunc(testStreamer))
	// Test
	c.AddFunc("@every 5m", recordFunc("u1_8ra"))

	c.Start()

	<-interrupt
	log.Println("Stopping")
}

func recordFunc(streamer string) func() {
	return func() {
		streamUrl, err := getStreamUrl(streamer)
		if err != nil {
			log.Printf("Error fetching stream URL for streamer [%s]: %v\n", streamer, err)
			return
		}
		log.Printf("Fetched stream URL for streamer [%s]: %s\n", streamer, streamUrl)

		sinkChan, err := makeFileSink(streamer)
		if err != nil {
			log.Println("Error creating recording file: ", err)
			return
		}

		record(streamer, streamUrl, sinkChan)
	}
}
