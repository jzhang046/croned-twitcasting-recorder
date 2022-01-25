package twitcasting

import (
	"context"
	"fmt"
	"log"

	"github.com/sacOO7/gowebsocket"

	"github.com/jzhang046/croned-twitcasting-recorder/record"
)

func RecordWS(
	recordContext record.RecordContext,
	cancelRecord context.CancelFunc,
	sinkChan chan<- []byte,
) {
	socket := gowebsocket.New(recordContext.GetStreamUrl())
	streamer := recordContext.GetStreamer()

	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		// Proxy: gowebsocket.BuildProxy("http://example.com"),
		UseSSL:         true,
		UseCompression: false,
		// Subprotocols:   []string{"chat", "superchat"},
	}

	socket.RequestHeader.Set("Origin", fmt.Sprintf("https://twitcasting.tv/%s", streamer))
	socket.RequestHeader.Set("User-Agent", userAgent)

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
