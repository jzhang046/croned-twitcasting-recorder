package main

import (
	"log"
)

const (
	testStreamer = "kaguramea_vov"
)

func main() {
	recordFunc(testStreamer)()
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
