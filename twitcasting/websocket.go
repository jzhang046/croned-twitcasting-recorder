package twitcasting

import (
	"log"

	"github.com/sacOO7/gowebsocket"

	"github.com/jzhang046/croned-twitcasting-recorder/record"
)

func RecordWS(recordCtx record.RecordContext, sinkChan chan<- []byte) {
	socket := gowebsocket.New(recordCtx.GetStreamUrl())
	streamer := recordCtx.GetStreamer()

	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		// Proxy: gowebsocket.BuildProxy("http://example.com"),
		UseSSL:         true,
		UseCompression: false,
		// Subprotocols:   []string{"chat", "superchat"},
	}

	socket.RequestHeader.Set("Origin", baseDomain)
	socket.RequestHeader.Set("User-Agent", userAgent)

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Error connecting to stream URL: ", err)
		recordCtx.Cancel()
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Printf("Connected to live stream for [%s], recording start \n", streamer)
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message", message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		sinkChan <- data
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Printf("Disconnected from live stream of [%s] \n", streamer)
		recordCtx.Cancel()
	}

	socket.Connect()

	// Waiting for context to finish
	<-recordCtx.Done()

	if socket.IsConnected {
		socket.Close()
	}
	close(sinkChan)
}
