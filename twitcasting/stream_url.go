package twitcasting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jmoiron/jsonq"
)

const (
	apiEndpoint    = "https://twitcasting.tv/streamserver.php"
	requestTimeout = 4 * time.Second
	userAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36"
)

var httpClient = &http.Client{
	Timeout: requestTimeout,
}

func GetWSStreamUrl(streamer string) (string, error) {
	u, _ := url.Parse(apiEndpoint)
	q := u.Query()
	q.Set("target", streamer)
	q.Set("mode", "client")
	u.RawQuery = q.Encode()

	request, _ := http.NewRequest("GET", u.String(), nil)
	request.UserAgent()
	request.Header.Set("User-Agent", userAgent)
	response, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("requesting stream info failed: %w", err)
	}
	defer response.Body.Close()

	responseData := map[string]interface{}{}
	dec := json.NewDecoder(response.Body)
	dec.Decode(&responseData)
	jq := jsonq.NewQuery(responseData)

	if err := checkStreamOnline(jq); err != nil {
		return "", err
	}

	// Try to get URL directly
	if streamUrl, err := getDirectStreamUrl(jq); err == nil {
		return streamUrl, nil
	}

	log.Printf("Direct Stream URL for streamer [%s] not available in the API response; fallback to default URL\n", streamer)
	return fallbackStreamUrl(jq, streamer)
}

func checkStreamOnline(jq *jsonq.JsonQuery) error {
	isLive, err := jq.Bool("movie", "live")
	if err != nil {
		return fmt.Errorf("error checking stream online status: %w", err)
	} else if !isLive {
		return fmt.Errorf("live stream is offline")
	}
	return nil
}

func getDirectStreamUrl(jq *jsonq.JsonQuery) (string, error) {
	// Try to get URL directly
	if streamUrl, err := jq.String("llfmp4", "streams", "main"); err == nil {
		return streamUrl, nil
	}
	if streamUrl, err := jq.String("llfmp4", "streams", "mobilesource"); err == nil {
		return streamUrl, nil
	}
	if streamUrl, err := jq.String("llfmp4", "streams", "base"); err == nil {
		return streamUrl, nil
	}

	return "", fmt.Errorf("direct stream URL not available")
}

func fallbackStreamUrl(jq *jsonq.JsonQuery, streamer string) (string, error) {
	mode := "base" // default mode
	if isSource, err := jq.Bool("fmp4", "source"); err == nil && isSource {
		mode = "main"
	} else if isMobile, err := jq.Bool("fmp4", "mobilesource"); err == nil && isMobile {
		mode = "mobilesource"
	}

	protocal, err := jq.String("fmp4", "proto")
	if err != nil {
		return "", fmt.Errorf("failed parsing stream protocal: %w", err)
	}

	host, err := jq.String("fmp4", "host")
	if err != nil {
		return "", fmt.Errorf("failed parsing stream host: %w", err)
	}

	movieId, err := jq.String("movie", "id")
	if err != nil {
		return "", fmt.Errorf("failed parsing movie ID: %w", err)
	}

	return fmt.Sprintf("%s:%s/ws.app/stream/%s/fmp4/bd/1/1500?mode=%s", protocal, host, movieId, mode), nil
}
