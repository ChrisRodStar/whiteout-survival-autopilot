package device

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/batazor/whiteout-survival-autopilot/internal/vision"
)

func (d *Device) handleEntryScreens(ctx context.Context) error {
	d.Logger.Info("🔎 Waiting for entry screens and pop-ups…")

	keywords := []string{
		"Welcome", "Alliance", "natalia", "Exploration", "Hero Gear",
		"General Speedup", "Construction Speedup", "Resource",
		"Mastery Material", "Purchase limit", "Agility",
		"Brothers in Arms", "Event Coming Soon", "Dawn Pack",
		"Unyielding Dawn", "Overview", "Confirm",
	}
	// Convert all keys to lowercase
	lowerKW := make([]string, len(keywords))
	for i, kw := range keywords {
		lowerKW[i] = strings.ToLower(kw)
	}

	start := time.Now()
	for {
		// Check overall timeout of 30 seconds
		if time.Since(start) > 30*time.Second {
			d.Logger.Info("⏱️ 30s elapsed, exiting without clicks")
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		zones, err := d.OCRClient.WaitForText(lowerKW, 10*time.Second, 1*time.Second, "entry_check")
		if err != nil {
			d.Logger.Error("❌ OCRClient error", slog.Any("err", err))
			return err
		}
		if len(zones) == 0 {
			continue
		}

		// 1) Look for Confirm
		for _, z := range zones {
			txt := strings.ToLower(strings.TrimSpace(z.Text))
			if vision.FuzzySubstringMatch(txt, "confirm", 1) &&
				z.AvgColor == "white" && z.BgColor == "green" {
				d.Logger.Info("🟢 Clicking Confirm", slog.String("text", txt))
				if err := d.ADB.ClickRegion("welcome_back_continue_button", d.AreaLookup); err != nil {
					d.Logger.Error("❌ Error clicking Confirm", slog.Any("err", err))
					return err
				}
				time.Sleep(time.Second)
				return nil
			}
		}

		// 2) Look for first pop-up
		found := false
		for _, z := range zones {
			txt := strings.ToLower(strings.TrimSpace(z.Text))
			for _, target := range lowerKW {
				if target == "confirm" {
					continue
				}
				if vision.FuzzySubstringMatch(txt, target, 1) {
					d.Logger.Info("🌀 Closing pop-up", slog.String("popup", txt))
					if err := d.ADB.ClickRegion("ad_banner_close", d.AreaLookup); err != nil {
						d.Logger.Error("❌ Error clicking pop-up close", slog.Any("err", err))
						return err
					}
					time.Sleep(300 * time.Millisecond)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		// If neither Confirm nor pop-up triggered — keep waiting
	}
}
