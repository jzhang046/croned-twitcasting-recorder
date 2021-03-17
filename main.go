package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sacOO7/gowebsocket"
)

const (
	user = "kaguramea_vov"

	timeFormt = "20060102-1504"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	streamUrl := getStreamUrl(user)
	log.Printf("Fetched stream URL: %s\n", streamUrl)
	socket := gowebsocket.New(streamUrl)

	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		// Proxy: gowebsocket.BuildProxy("http://example.com"),
		UseSSL:         true,
		UseCompression: false,
		// Subprotocols:   []string{"chat", "superchat"},
	}

	socket.RequestHeader.Set("Origin", fmt.Sprintf("https://twitcasting.tv/%s", user))
	socket.RequestHeader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Recieved connect error ", err)
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message  " + message)
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(fmt.Sprintf("%s-%s.ts", user, time.Now().Format(timeFormt)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		fmt.Print(".")
		if _, err := f.Write(data); err != nil {
			log.Fatal(err)
		}
	}
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}
	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
