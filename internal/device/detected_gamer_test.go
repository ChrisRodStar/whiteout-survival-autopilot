package device_test

import (
	"context"
	"image"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/device"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/fsm"
)

type MockADB struct {
	imagePath string
}

func (m *MockADB) ClickRegion(name string, lookup *config.AreaLookup) error {
	return nil
}

func (m *MockADB) ClickOCRResult(_ *domain.OCRResult) error {
	return nil
}

func (m *MockADB) Swipe(x1 int, y1 int, x2 int, y2 int, durationMs time.Duration) error {
	return nil
}

func (m *MockADB) ListDevices() ([]string, error) {
	return nil, nil
}

func (m *MockADB) SetActiveDevice(serial string) {}

func (m *MockADB) GetActiveDevice() string {
	return ""
}

func TestDetectedGamer_WithRealConfig_AndDeviceNew(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// 🧾 Load config
	cfg, err := config.LoadDeviceConfig("../../db/devices.yaml", "../../db/state.yaml")
	if err != nil {
		t.Fatalf("❌ Failed to load devices.yaml: %v", err)
	}

	if len(cfg.Devices) == 0 || len(cfg.Devices[0].Profiles) == 0 {
		t.Fatal("❌ No devices or profiles in config")
	}

	// ⚙️ Replace ADB controller with mock
	profiles := cfg.Devices[0].Profiles
	log := slog.Default()

	// Create Device via `New`, then replace ADB and FSM
	dev, err := device.New("test-device", profiles, log, "../../references/area.json")
	if err != nil {
		t.Fatalf("❌ device.New() returned error: %v", err)
	}

	// Load area.json
	lookup, err := config.LoadAreaReferences("../../references/area.json")
	if err != nil {
		t.Fatalf("failed to load area.json: %v", err)
	}

	dev.FSM = fsm.NewGame(log, &MockADB{imagePath: "../../references/screenshots/chief_profile.png"}, lookup)

	// Replace only Screenshot logic, the rest can remain
	dev.ADB = &MockADB{imagePath: "../../references/screenshots/chief_profile.png"}

	// 🚀 Perform detection
	profileIdx, gamerIdx, err := dev.DetectedGamer(ctx, "../../references/screenshots/chief_profile.png")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, profileIdx, 0, "should find profile")
	assert.GreaterOrEqual(t, gamerIdx, 0, "should find gamer")

	nickname := dev.Profiles[profileIdx].Gamer[gamerIdx].Nickname
	t.Logf("✅ Player found: profileIdx=%d, gamerIdx=%d, nickname=%s", profileIdx, gamerIdx, nickname)

	assert.Equal(t, "batazor", nickname)
}
