package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/jzhang046/croned-twitcasting-recorder/config"
	"github.com/jzhang046/croned-twitcasting-recorder/record"
	"github.com/jzhang046/croned-twitcasting-recorder/sink"
	"github.com/jzhang046/croned-twitcasting-recorder/twitcasting"
)

const terminationGracePeriodSeconds = 3

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	log.Println("croned recorder starting ")

	config := config.GetDefaultConfig()
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	rootCtx, cancalAllRecords := context.WithCancel(context.Background())

	for _, streamerConfig := range config.Streamers {
		if _, err := c.AddFunc(
			streamerConfig.Schedule,
			record.ToRecordFunc(&record.RecordConfig{
				Streamer:         streamerConfig.ScreenId,
				StreamUrlFetcher: twitcasting.GetWSStreamUrl,
				SinkProvider:     sink.NewFileSink,
				StreamRecorder:   twitcasting.RecordWS,
				RootContext:      rootCtx,
			}),
		); err != nil {
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

	log.Printf("Terminating in %d seconds.. \n", terminationGracePeriodSeconds)
	go func() {
		cancalAllRecords()
		c.Stop()
	}()

	time.Sleep(terminationGracePeriodSeconds * time.Second)
	log.Fatal("Terminated")
}
