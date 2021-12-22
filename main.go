package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"

	"github.com/jzhang046/croned-twitcasting-recorder/config"
	"github.com/jzhang046/croned-twitcasting-recorder/record"
	"github.com/jzhang046/croned-twitcasting-recorder/sink"
	"github.com/jzhang046/croned-twitcasting-recorder/twitcasting"
)

func main() {
	log.Println("croned recorder starting ")

	config := config.GetDefaultConfig()
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	for _, streamerConfig := range config.Streamers {
		if _, err := c.AddFunc(streamerConfig.Schedule, recordFunc(streamerConfig.ScreenId)); err != nil {
			log.Fatalln("Failed adding record schedule: ", err)
		} else {
			log.Printf("Added schedule [%s] for streamer [%s] \n", streamerConfig.Schedule, streamerConfig.ScreenId)
		}
	}

	c.Start()
	log.Println("croned recorder started ")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	log.Fatal("Terminated")
}

func recordFunc(streamer string) func() {
	return func() {
		streamUrl, err := twitcasting.GetWSStreamUrl(streamer)
		if err != nil {
			log.Printf("Error fetching stream URL for streamer [%s]: %v\n", streamer, err)
			return
		}
		log.Printf("Fetched stream URL for streamer [%s]: %s\n", streamer, streamUrl)
		recordContext, cancelRecord := record.NewRecordContext(streamer, streamUrl)

		sinkChan, err := sink.NewFileSink(streamer, cancelRecord)
		if err != nil {
			log.Println("Error creating recording file: ", err)
			return
		}

		twitcasting.RecordWS(recordContext, cancelRecord, sinkChan)
	}
}
