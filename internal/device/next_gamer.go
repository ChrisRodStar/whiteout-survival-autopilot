package device

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
	"github.com/batazor/whiteout-survival-autopilot/internal/fsm"
)

func (d *Device) NextGamer(profileIdx, gamerIdx int) {
	// Initialize context and span
	ctx := context.Background()
	tracer := otel.Tracer("device")
	ctx, span := tracer.Start(ctx, "NextGamer")
	defer span.End()

	// Extract traceID for logs
	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

	d.activeProfileIdx = profileIdx
	d.activeGamerIdx = gamerIdx

	profile := d.Profiles[profileIdx]
	gamer := &profile.Gamer[gamerIdx]

	d.Logger.Info("🎮 Switching to another player in current profile",
		slog.String("email", profile.Email),
		slog.String("nickname", gamer.Nickname),
		slog.Int("id", gamer.ID),
		slog.String("trace_id", traceID),
	)

	// Set new player in FSM
	d.FSM.SetCallback(gamer)

	// 🔁 Navigation: go to character selection screen
	d.Logger.Info("➡️ Navigating to player selection screen",
		slog.String("trace_id", traceID),
	)
	d.FSM.ForceTo(state.StateChiefCharacters, nil)

	// 🕒 Wait to avoid conflicts with other processes
	time.Sleep(2 * time.Second)

	// ========== 1️⃣ Perform unified full-screen OCR ==========
	fullOCR, fullErr := d.OCRClient.FetchOCR("", nil) // debugName can be omitted
	if fullErr != nil {
		d.Logger.Error("❌ Full OCR failed", slog.Any("error", fullErr))
		panic(fmt.Sprintf("ocrClient.FetchOCR() failed: %v", fullErr))
	}

	// Wait for nickname
	var gamerZone *domain.OCRResult
	for _, zone := range fullOCR {
		if strings.Contains(zone.Text, gamer.Nickname) {
			gamerZone = &zone
			break
		}
	}

	d.Logger.Info("🟢 Clicking player nickname",
		slog.String("text", gamerZone.Text),
		slog.String("trace_id", traceID),
	)
	if err := d.ADB.ClickOCRResult(gamerZone); err != nil {
		d.Logger.Error("❌ Failed to click nickname account",
			slog.Any("err", err),
			slog.String("trace_id", traceID),
		)
		panic(fmt.Sprintf("ClickRegion(nickname:%s) failed: %v", gamer.Nickname, err))
	}

	time.Sleep(2 * time.Second)

	d.Logger.Info("🟢 Clicking confirmation button",
		slog.String("region", "character_change_confirm"),
		slog.String("trace_id", traceID),
	)
	if err := d.ADB.ClickRegion("character_change_confirm", d.AreaLookup); err != nil {
		d.Logger.Error("❌ Failed to click character_change_confirm",
			slog.Any("err", err),
			slog.String("trace_id", traceID),
		)
		panic(fmt.Sprintf("ClickRegion(character_change_confirm) failed: %v", err))
	}

	// Check entry banners
	err := d.handleEntryScreens(ctx)
	if err != nil {
		d.Logger.Error("❌ Failed to handle entry banners",
			slog.Any("err", err),
			slog.String("trace_id", traceID),
		)
		panic(fmt.Sprintf("handleEntryScreens() failed: %v", err))
	}

	d.Logger.Info("✅ Login completed, navigating to Main City",
		slog.String("trace_id", traceID),
	)
	d.Logger.Info("🔧 Initializing FSM",
		slog.String("trace_id", traceID),
	)
	d.FSM = fsm.NewGame(d.Logger, d.ADB, d.AreaLookup, d.triggerEvaluator, d.ActiveGamer(), d.OCRClient)
}
