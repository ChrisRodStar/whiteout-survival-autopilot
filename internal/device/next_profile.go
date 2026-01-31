package device

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
	"github.com/batazor/whiteout-survival-autopilot/internal/fsm"
	"github.com/batazor/whiteout-survival-autopilot/internal/ocrclient"
)

func (d *Device) NextProfile(profileIdx, expectedGamerIdx int) {
	// 🕒 Wait to avoid conflicts with other processes
	time.Sleep(500 * time.Millisecond)

	ctx := context.Background()

	profile := d.Profiles[profileIdx]
	expected := &profile.Gamer[expectedGamerIdx]

	d.Logger.Info("🎮 Changing active player",
		slog.String("email", profile.Email),
		slog.String("expected", expected.Nickname),
	)

	// 🔁 Navigation: go to Google account selection screen
	d.Logger.Info("➡️ Navigating to account selection screen")
	d.FSM.ForceTo(state.StateChiefProfileAccountChangeGoogle, nil)

	// 🕒 Wait to avoid conflicts with other processes
	time.Sleep(2 * time.Second)

	// ========== 1️⃣ Perform unified full-screen OCR ==========
	region, ok := d.AreaLookup.Get("google_profile")
	if !ok {
		d.Logger.Error("❌ Failed to find google_profile area")
		panic("AreaLookup(google_profile) failed")
	}

	fullOCR, fullErr := d.OCRClient.FetchOCR("google_profile", []ocrclient.Region{
		{
			X0: region.Zone.Min.X,
			Y0: region.Zone.Min.Y,
			X1: region.Zone.Max.X,
			Y1: region.Zone.Max.Y,
		},
	})
	if fullErr != nil {
		d.Logger.Error("❌ Full OCR failed", slog.Any("error", fullErr))
		panic(fmt.Sprintf("ocrClient.FetchOCR() failed: %v", fullErr))
	}

	// 📦 OCR by email
	var emailZone *domain.OCRResult
	for _, zone := range fullOCR {
		if zone.Text == profile.Email {
			emailZone = &zone
			break
		}
	}

	d.Logger.Info("🟢 Clicking email account", slog.String("text", emailZone.Text), slog.String("region", emailZone.String()))
	if err := d.ADB.ClickOCRResult(emailZone); err != nil {
		d.Logger.Error("❌ Failed to click email account", slog.Any("err", err))
		panic(fmt.Sprintf("ClickRegion(email:gamer1) failed: %v", err))
	}

	time.Sleep(3 * time.Second)

	googleContinueArea, ok := d.AreaLookup.Get("to_google_continue")
	if !ok {
		d.Logger.Error("❌ Failed to find to_google_continue area")
		panic("AreaLookup(to_google_continue) failed")
	}

	// Wait for "Continue" text via OCR client
	if _, err := d.OCRClient.WaitForText([]string{"Continue"}, time.Second, 500*time.Millisecond, "continue"); err != nil {
		d.Logger.Error("❌ OCRClient WaitForText failed for Continue", slog.Any("err", err))
		panic(fmt.Sprintf("OCRClient.WaitForText(Continue) failed: %v", err))
	}

	d.Logger.Info("🟢 Clicking Google continue button", slog.String("region", "to_google_continue"))

	if err := d.ADB.Click(googleContinueArea.Zone); err != nil {
		d.Logger.Error("❌ Failed to click to_google_continue", slog.Any("err", err))
		panic(fmt.Sprintf("ClickRegion(to_google_continue) failed: %v", err))
	}

	// ♻️ Reset FSM after login
	d.activeProfileIdx = profileIdx
	d.activeGamerIdx = expectedGamerIdx
	d.FSM = fsm.NewGame(d.Logger, d.ADB, d.AreaLookup, d.triggerEvaluator, d.ActiveGamer(), d.OCRClient)

	// Check entry banners
	err := d.handleEntryScreens(ctx)
	if err != nil {
		d.Logger.Error("❌ Failed to handle entry banners", slog.Any("err", err))
		panic(fmt.Sprintf("handleEntryScreens() failed: %v", err))
	}

	// 🔍 Check that the active profile is the one expected
	active, pIdx, _, err := d.DetectAndSetCurrentGamer(ctx)
	if err != nil || pIdx != profileIdx {
		d.Logger.Warn("⚠️ After login, active profile doesn't match", slog.Any("detected_profile", pIdx), slog.Any("err", err))
		return
	}

	// 🧾 If player is not the right one — switch manually
	if active.ID != expected.ID {
		d.Logger.Warn("🛑 Automatically selected wrong player — switching",
			slog.String("expected", expected.Nickname),
			slog.String("got", active.Nickname),
		)
		d.NextGamer(profileIdx, expectedGamerIdx)
	}

	// ✅ Set callback
	d.FSM.SetCallback(active)

	d.Logger.Info("✅ Successfully switched to new profile", "nickname", active.Nickname)
}

func (d *Device) ActiveGamer() *domain.Gamer {
	if d.activeProfileIdx >= 0 && d.activeProfileIdx < len(d.Profiles) {
		profile := d.Profiles[d.activeProfileIdx]
		if d.activeGamerIdx >= 0 && d.activeGamerIdx < len(profile.Gamer) {
			return &profile.Gamer[d.activeGamerIdx]
		}
	}
	return nil
}
