package cmd

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/jzhang046/croned-twitcasting-recorder/config"
	"github.com/jzhang046/croned-twitcasting-recorder/record"
	"github.com/jzhang046/croned-twitcasting-recorder/sink"
	"github.com/jzhang046/croned-twitcasting-recorder/twitcasting"
)

const CronedRecordCmdName = "croned"

func RecordCroned() {
	log.Printf("Starting in recoding mode [%s].. \n", CronedRecordCmdName)

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

	<-waitForInterruput(func() {
		cancalAllRecords()
		<-c.Stop().Done()
	})

	log.Fatal("Terminated on user interrupt")
}
