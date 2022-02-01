package cmd

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"

	"github.com/jzhang046/croned-twitcasting-recorder/config"
	"github.com/jzhang046/croned-twitcasting-recorder/record"
	"github.com/jzhang046/croned-twitcasting-recorder/sink"
	"github.com/jzhang046/croned-twitcasting-recorder/twitcasting"
)

const CronedRecordCmdName = "croned"

func RecordCroned() {
	log.Printf("Starting in recoding mode [%s] with PID [%d].. \n", CronedRecordCmdName, os.Getpid())

	config := config.GetDefaultConfig()
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	interruptCtx, afterGracefulInterrupt := newInterruptableCtx()

	for _, streamerConfig := range config.Streamers {
		if _, err := c.AddFunc(
			streamerConfig.Schedule,
			record.ToRecordFunc(&record.RecordConfig{
				Streamer:         streamerConfig.ScreenId,
				StreamUrlFetcher: twitcasting.GetWSStreamUrl,
				SinkProvider:     sink.NewFileSink,
				StreamRecorder:   twitcasting.RecordWS,
				RootContext:      interruptCtx,
			}),
		); err != nil {
			log.Fatalln("Failed adding record schedule: ", err)
		} else {
			log.Printf("Added schedule [%s] for streamer [%s] \n", streamerConfig.Schedule, streamerConfig.ScreenId)
		}
	}

	c.Start()
	log.Println("croned recorder started ")

	// interrupt => stop cron and wait for all task to complete => wait for graceful interrupt
	<-interruptCtx.Done()
	<-c.Stop().Done()
	<-afterGracefulInterrupt

	log.Fatal("Terminated on user interrupt")
}
