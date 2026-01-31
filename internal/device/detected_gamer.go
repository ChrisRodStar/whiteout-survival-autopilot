package device

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
	"github.com/batazor/whiteout-survival-autopilot/internal/ocrclient"
)

func (d *Device) DetectedGamer(ctx context.Context) (int, int, error) {
	d.Logger.Info("🚀 Detecting current player")

	// 0. Navigate to profile screen
	d.FSM.ForceTo(state.StateChiefProfile, nil)

	defer func() {
		// 4. Return to main screen
		d.FSM.ForceTo(state.StateMainCity, nil)
	}()

	zone, ok := d.AreaLookup.Get("chief_profile_nickname")
	if !ok {
		d.Logger.Error("GetRegionByName failed",
			slog.String("region", "chief_profile_nickname"),
		)
		return -1, -1, errors.New("no matches found with nickname")
	}

	region := ocrclient.Region{
		X0: zone.Zone.Min.X,
		Y0: zone.Zone.Min.Y,
		X1: zone.Zone.Max.X,
		Y1: zone.Zone.Max.Y,
	}

	// 3. Recognize player nickname
	fullOCR, fullErr := d.OCRClient.FetchOCR("", []ocrclient.Region{region}) // debugName can be omitted
	if fullErr != nil {
		d.Logger.Error("Full OCR failed", slog.Any("error", fullErr))
		return -1, -1, fmt.Errorf("full OCR failed: %w", fullErr)
	}

	if len(fullOCR) == 0 {
		d.Logger.Warn("⚠️ Failed to recognize player nickname", slog.String("region", "chief_profile_nickname"))
		return -1, -1, errors.New("no matches found with nickname")
	}

	nicknameParsed := fullOCR[0].Text

	// drop aliance [RLX]batazor -> batazor
	if strings.Contains(nicknameParsed, "]") {
		nicknameParsed = strings.Split(nicknameParsed, "]")[1]
	}

	d.Logger.Info("🟢 Nickname recognized", slog.String("parsed", nicknameParsed))

	type matchInfo struct {
		profileIdx int
		gamerIdx   int
		score      int
	}

	var matches []matchInfo

	for pIdx, profile := range d.Profiles {
		for gIdx, gamer := range profile.Gamer {
			expected := strings.ToLower(strings.TrimSpace(gamer.Nickname))
			if matched := fuzzy.RankMatch(expected, nicknameParsed); matched != -1 {
				matches = append(matches, matchInfo{pIdx, gIdx, matched})
			}
		}
	}

	if len(matches) == 0 {
		d.Logger.Warn("⚠️ Nickname not found by fuzzy match", slog.String("parsed", nicknameParsed))
		return -1, -1, errors.New("no matches found with nickname")
	}

	// Find the best match (with the lowest score)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score < matches[j].score
	})
	best := matches[0]

	d.Logger.Info("✅ Player found",
		slog.Int("profileIdx", best.profileIdx),
		slog.Int("gamerIdx", best.gamerIdx),
		slog.Int("score", best.score),
	)

	return best.profileIdx, best.gamerIdx, nil
}

func (d *Device) DetectAndSetCurrentGamer(ctx context.Context) (*domain.Gamer, int, int, error) {
	pIdx, gIdx, err := d.DetectedGamer(ctx)
	if err != nil || pIdx < 0 || gIdx < 0 {
		d.Logger.Warn("⚠️ Failed to detect active player", slog.Any("err", err))
		return nil, -1, -1, err
	}

	// 💾 Save as current
	d.activeProfileIdx = pIdx
	d.activeGamerIdx = gIdx

	active := &d.Profiles[pIdx].Gamer[gIdx]
	d.Logger.Info("🔎 Active player detected", slog.String("nickname", active.Nickname))

	d.FSM.SetCallback(active)

	// Reset old state
	active.ScreenState.Reset()

	return active, pIdx, gIdx, nil
}
