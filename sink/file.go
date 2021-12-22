package sink

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	timeFormt      = "20060102-1504"
	sinkChanBuffer = 16
)

func NewFileSink(streamer string, cancelRecord context.CancelFunc) (chan []byte, error) {
	// If the file doesn't exist, create it, or append to the file
	filename := fmt.Sprintf("%s-%s.ts", streamer, time.Now().Format(timeFormt))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.Printf("Recording file %s", filename)

	sinkChan := make(chan []byte, sinkChanBuffer)

	go func() {
		defer f.Close()
		for {
			data, more := <-sinkChan
			if more {
				if _, err := f.Write(data); err != nil {
					log.Printf("Error writing recording file %s: %v\n", filename, err)
					cancelRecord()
					return
				}
			} else {
				log.Printf("Completed writing all data to %s\n", filename)
				return
			}
		}
	}()

	return sinkChan, nil
}
