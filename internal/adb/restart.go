package adb

import (
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

const gamePackageName = "com.gof.global"

// RestartApplication restarts the application via adb commands.
func (a *Controller) RestartApplication() error {
	a.logger.Warn("🔄 Restarting application", slog.String("package", gamePackageName))

	// Close the application
	closeCmd := exec.Command("adb", "-s", a.deviceID, "shell", "am", "force-stop", gamePackageName)
	if err := closeCmd.Run(); err != nil {
		a.logger.Error("❌ Error closing application", slog.String("package", gamePackageName), slog.Any("error", err))
		return fmt.Errorf("failed to close app %s: %w", gamePackageName, err)
	}

	time.Sleep(2 * time.Second)

	// Start the application again
	startCmd := exec.Command("adb", "-s", a.deviceID, "shell", "monkey", "-p", gamePackageName, "-c", "android.intent.category.LAUNCHER", "1")
	if err := startCmd.Run(); err != nil {
		a.logger.Error("❌ Error starting application", slog.String("package", gamePackageName), slog.Any("error", err))
		return fmt.Errorf("failed to start app %s: %w", gamePackageName, err)
	}

	a.logger.Info("✅ Application successfully restarted", slog.String("package", gamePackageName))

	return nil
}
