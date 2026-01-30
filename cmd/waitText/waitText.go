package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/batazor/whiteout-survival-autopilot/internal/ocrclient"
)

func main() {
	// 1) Initialize client
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	client := ocrclient.NewClient("RF8RC00M8MF", logger)

	// 2) Set stop words, timeout and polling interval
	stopWords := []string{"Chief Profile"}
	timeout := 30 * time.Second // wait maximum 30 seconds
	interval := 1 * time.Second // check screen every 1 second
	debugName := "screenState.titleFact"

	// 3) Start waiting
	results, err := client.WaitForText(stopWords, timeout, interval, debugName)
	if err != nil {
		logger.Error("WaitForText failed", "error", err)
		return
	}

	// 4) Process result — list of OCR zones where text was found
	if len(results) == 0 {
		fmt.Println("⚠️ None of the stop words found")
		return
	}

	fmt.Println("✅ Found one of the words! Zones with recognized text:")
	for i, res := range results {
		fmt.Printf("  %d) \"%s\" (%.2f) @ X:%d Y:%d W:%d H:%d\n",
			i+1, res.Text, res.Score, res.X, res.Y, res.Width, res.Height)
	}
}
