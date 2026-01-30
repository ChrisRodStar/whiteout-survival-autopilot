package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/batazor/whiteout-survival-autopilot/internal/ocrclient"
)

func main() {
	// 1) Create OCR client
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	client := ocrclient.NewClient("RF8RC00M8MF", logger)

	// 2) Try to find icon "alliance.state.isNeedSupport" (file name alliance.state.isNeedSupport.png in references/icons),
	//    with threshold 0.8 and debug_name label "alliance.state.isNeedSupport_check"
	start := time.Now()
	resp, err := client.FindImage("alliance.state.isNeedSupport", 0.8, "alliance.state.isNeedSupport_check")
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("FindImage failed: %v", err)
	}

	// 3) Parse result
	if resp.Found {
		// convert to image.Rectangle rectangles
		rects := resp.ToRects()
		fmt.Printf("✅ Found icon «alliance.state.isNeedSupport» (threshold=0.8) in %v:\n", elapsed)
		for i, r := range rects {
			fmt.Printf("  #%d at x=%d,y=%d – x2=%d,y2=%d\n",
				i+1, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
		}
	} else {
		fmt.Printf("❌ Icon «alliance.state.isNeedSupport» not found (checked in %v)\n", elapsed)
	}
}
