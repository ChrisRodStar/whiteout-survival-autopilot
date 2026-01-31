package device

//
//import (
//	"fmt"
//	"log/slog"
//	"strings"
//	"time"
//
//	"github.com/batazor/whiteout-survival-autopilot/internal/adb"
//	"github.com/batazor/whiteout-survival-autopilot/internal/config"
//	"github.com/batazor/whiteout-survival-autopilot/internal/ocrclient"
//	"github.com/batazor/whiteout-survival-autopilot/internal/vision"
//)
//
//type ReconnectHandler struct {
//	adbController adb.DeviceController
//	area          *config.AreaLookup
//	OCRClient     *ocrclient.Client
//	logger        *slog.Logger
//	maxAttempts   int
//}
//
//func (h *ReconnectHandler) HandleReconnect(screenshotPath string) error {
//	const waitAfterReconnectClick = 20 * time.Second
//	const waitAfterRestart = 10 * time.Second
//	const maxTimeout = 20 * time.Second
//
//	for restartCount := 0; ; restartCount++ {
//		attempt := 0
//
//		for attempt < h.maxAttempts {
//			attempt++
//
//			found, err := h.checkReconnectWindow(screenshotPath)
//			if err != nil {
//				h.logger.Warn("❌ Error checking reconnect window", slog.Any("err", err))
//				return err
//			}
//			if !found {
//				h.logger.Info("✅ Reconnect button not detected, continuing work")
//				return nil
//			}
//
//			h.logger.Warn("Reconnect window detected, trying to reconnect",
//				slog.Int("attempt", attempt),
//				slog.Int("restartCount", restartCount),
//			)
//
//			if err := h.adbController.ClickRegion("reconnect_button", h.area); err != nil {
//				return fmt.Errorf("error clicking reconnect button: %w", err)
//			}
//
//			h.logger.Info("⏳ Waiting for loading to complete after click (20 seconds)")
//			time.Sleep(waitAfterReconnectClick)
//		}
//
//		// --- Application restart ---
//		h.logger.Error("🚨 Reconnection failed, restarting application")
//		if err := h.adbController.RestartApplication(); err != nil {
//			return fmt.Errorf("failed to restart application: %w", err)
//		}
//
//		h.logger.Info("⏳ Waiting for application to load (10 seconds)")
//		time.Sleep(waitAfterRestart)
//
//		// --- After restart, wait for reconnect button to disappear ---
//		h.logger.Info("⏳ Waiting for reconnect button to disappear (up to 20 seconds)")
//
//		expire := time.After(maxTimeout)
//		tick := time.NewTicker(2 * time.Second)
//		defer tick.Stop()
//
//		for {
//			select {
//			case <-expire:
//				h.logger.Warn("🔁 Reconnect button still on screen after restart — continuing loop")
//				break
//
//			case <-tick.C:
//				found, err := h.checkReconnectWindow(screenshotPath)
//				if err != nil {
//					return err
//				}
//				if !found {
//					h.logger.Info("✅ Reconnect button disappeared — continuing work")
//					return nil
//				}
//			}
//		}
//	}
//}
//
//func (h *ReconnectHandler) checkReconnectWindow(screenshotPath string) (bool, error) {
//	// perform OCR only in the reconnect button area
//	results, err := h.OCRClient.FetchOCRByAreaName("reconnect_button", "reconnect_check")
//	if err != nil {
//		h.logger.Error("❌ OCRClient FetchOCRByAreaName failed for reconnect", slog.Any("error", err))
//		return false, err
//	}
//
//	// if nothing recognized — window didn't appear
//	if len(results) == 0 {
//		return false, nil
//	}
//
//	// take first result and convert to lowercase
//	text := strings.ToLower(strings.TrimSpace(results[0].Text))
//	h.logger.Info("🔍 OCR result reconnect", slog.String("text", text))
//
//	// fuzzy match for the word "reconnect"
//	target := "reconnect"
//	if strings.Contains(text, target) || vision.FuzzySubstringMatch(text, target, 1) {
//		h.logger.Info("✅ reconnect window detected")
//		return true, nil
//	}
//
//	return false, nil
//}
