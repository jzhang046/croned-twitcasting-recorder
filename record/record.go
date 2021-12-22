package record

import (
	"context"
	"log"
)

type RecordConfig struct {
	Streamer         string
	StreamUrlFetcher func(string) (string, error)
	SinkProvider     func(string, context.CancelFunc) (chan<- []byte, error)
	StreamRecorder   func(RecordContext, context.CancelFunc, chan<- []byte)
	RootContext      context.Context
}

func ToRecordFunc(recordConfig *RecordConfig) func() {
	streamer := recordConfig.Streamer
	return func() {
		streamUrl, err := recordConfig.StreamUrlFetcher(streamer)
		if err != nil {
			log.Printf("Error fetching stream URL for streamer [%s]: %v\n", streamer, err)
			return
		}
		log.Printf("Fetched stream URL for streamer [%s]: %s\n", streamer, streamUrl)
		recordContext, cancelRecord := newRecordContext(recordConfig.RootContext, streamer, streamUrl)

		sinkChan, err := recordConfig.SinkProvider(streamer, cancelRecord)
		if err != nil {
			log.Println("Error creating recording file: ", err)
			return
		}

		recordConfig.StreamRecorder(recordContext, cancelRecord, sinkChan)
	}
}
