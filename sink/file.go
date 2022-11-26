package sink

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jzhang046/croned-twitcasting-recorder/record"
)

const (
	timeFormat     = "20060102-1504"
	sinkChanBuffer = 16
)

func NewFileSink(recordCtx record.RecordContext) (chan<- []byte, error) {
	// If the file doesn't exist, create it, or append to the file
	filename := fmt.Sprintf("%s-%s.ts", recordCtx.GetStreamer(), time.Now().Format(timeFormat))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	log.Printf("Recording file %s", filename)

	sinkChan := make(chan []byte, sinkChanBuffer)

	go func() {
		defer f.Close()
		for data := range sinkChan {
			if _, err = f.Write(data); err != nil {
				log.Printf("Error writing recording file %s: %v\n", filename, err)
				recordCtx.Cancel()
				return
			}
		}
		log.Printf("Completed writing all data to %s\n", filename)
	}()

	return sinkChan, nil
}
