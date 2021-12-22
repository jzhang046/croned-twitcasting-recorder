package twitcasting

import (
	"context"
	"fmt"
	"log"

	"github.com/sacOO7/gowebsocket"
)

func RecordWS(streamer, streamUrl string, sinkChan chan<- []byte) {
	socket := gowebsocket.New(streamUrl)

	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		// Proxy: gowebsocket.BuildProxy("http://example.com"),
		UseSSL:         true,
		UseCompression: false,
		// Subprotocols:   []string{"chat", "superchat"},
	}

	socket.RequestHeader.Set("Origin", fmt.Sprintf("https://twitcasting.tv/%s", streamer))
	socket.RequestHeader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	recordContext, cancelRecord := context.WithCancel(context.Background())

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Error connecting to stream URL: ", err)
		cancelRecord()
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Printf("Connected to live stream for [%s], recording start \n", streamer)
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Unable to continue recording for [%s]: %s \n", streamer, r)
				cancelRecord()
			}
		}()
		sinkChan <- data
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Printf("Disconnected from live stream of [%s] \n", streamer)
		cancelRecord()
	}

	socket.Connect()

	// Waiting for context to finish
	<-recordContext.Done()

	if socket.IsConnected {
		socket.Close()
	}
	close(sinkChan)
}
