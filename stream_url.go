package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/jmoiron/jsonq"
)

const (
	apiEndpoint = "https://twitcasting.tv/streamserver.php"
)

func getStreamUrl(streamer string) string {
	u, _ := url.Parse(apiEndpoint)
	q := u.Query()
	q.Set("target", streamer)
	q.Set("mode", "client")
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		log.Fatal("HTTP request failed", err)
	}

	responseData := map[string]interface{}{}
	dec := json.NewDecoder(response.Body)
	dec.Decode(&responseData)
	jq := jsonq.NewQuery(responseData)

	isLive, err := jq.Bool("movie", "live")
	if err != nil || !isLive {
		log.Fatalf("Live stream of %s is offline", streamer)
	}

	// Try to get URL directly
	if streamUrl, err := jq.String("llfmp4", "streams", "main"); err == nil {
		return streamUrl
	}
	if streamUrl, err := jq.String("llfmp4", "streams", "mobilesource"); err == nil {
		return streamUrl
	}
	if streamUrl, err := jq.String("llfmp4", "streams", "base"); err == nil {
		return streamUrl
	}

	log.Println("Stream URL not directly available in the API response; fallback to default URL")
	mode := "base" // default mode
	if isSource, err := jq.Bool("fmp4", "source"); err == nil && isSource {
		mode = "main"
	} else if isMobile, err := jq.Bool("fmp4", "mobilesource"); err == nil && isMobile {
		mode = "mobilesource"
	}

	protocal, err := jq.String("fmp4", "proto")
	if err != nil {
		log.Fatal("Failed to parse protocal", err)
	}
	host, err := jq.String("fmp4", "host")
	if err != nil {
		log.Fatal("Failed to parse host", err)
	}
	movieId, err := jq.String("movie", "id")
	if err != nil {
		log.Fatal("Failed to parse movie ID", err)
	}
	return fmt.Sprintf("%s:%s/ws.app/stream/%s/fmp4/bd/1/1500?mode=%s", protocal, host, movieId, mode)
}
