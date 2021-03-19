package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
)

func main() {
	log.Println("croned recorder started ")

	config := getDefaultConfig()
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	for _, sc := range config.Streamers {
		if _, err := c.AddFunc(sc.Schedule, recordFunc(sc.ScreenId)); err != nil {
			log.Fatalln("Failed adding record schedule: ", err)
		} else {
			log.Printf("Added schedule [%s] for streamer [%s] \n", sc.Schedule, sc.ScreenId)
		}
	}

	c.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt
	log.Fatal("Terminated")
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
