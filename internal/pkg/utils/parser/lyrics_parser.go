package parser

import (
	"encoding/json"
	"strings"
)

func ParseLyrics(rawLyrics string) ([]byte, error) {
	lines := strings.Split(strings.ReplaceAll(rawLyrics, "\r\n", "\n"), "\n")

	var verses []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			verses = append(verses, trimmedLine)
		}
	}

	lyricsData := struct {
		Text   string   `json:"text"`
		Verses []string `json:"verses"`
	}{
		Text:   rawLyrics,
		Verses: verses,
	}

	return json.Marshal(lyricsData)
}
