package twitcasting

/**
Example streamserver.php API response:
{
  "movie": {
    "id": 1234, // number
    "live": true
  },
  "hls": {
    "host": "twitcasting.tv",
    "proto": "https",
    "source": false
  },
  "fmp4": {
    "host": "10-0-0-1.twitcasting.tv", // Decimal IP number separated by dash
    "proto": "wss",
    "source": false,
    "mobilesource": false
  },
  "llfmp4": {
    "streams": {
      "main": "wss://10-0-0-1.twitcasting.tv/tc.edge/v1/streams/1234.567.89/fmp4"
    }
  }
}
*/
