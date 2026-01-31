package fsm

import (
	"log/slog"

	"github.com/samber/lo"

	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
	"github.com/batazor/whiteout-survival-autopilot/internal/vision"
)

// ExpectState checks the current screen state.
// If the title matches a known screen but the state doesn't match — returns the expected one.
func (g *GameFSM) ExpectState(want string) (string, error) {
	gamerState, err := g.analyzer.AnalyzeAndUpdateState(
		g.gamerState, g.rulesCheckState["default"], nil,
	)
	if err != nil {
		g.logger.Error("❌ Failed to analyze screen",
			slog.String("action", "expect_state"),
			slog.String("state", want),
			slog.Any("error", err),
		)
		return "", err
	}

	ocrTitle := gamerState.ScreenState.TitleFact
	ocrFamily := gamerState.ScreenState.IsMainCity

	g.logger.Debug("FSM: state analysis",
		slog.String("ocr_title", ocrTitle),
		slog.String("ocr_family", ocrFamily),
		slog.String("want", want),
	)

	// 0. If family is defined — use only its group
	switch {
	case vision.FuzzySubstringMatch(ocrFamily, "world", 1):
		g.logger.Debug("FSM: Family defined as CITY (by world reference)",
			slog.String("ocr_family", ocrFamily),
			slog.String("want", want),
		)
		groupStates, ok := config.TitleToState["MainCity"]
		g.logger.Debug("FSM: Group MainCity", slog.Any("groupStates", groupStates), slog.Bool("found", ok))
		if ok && lo.Contains(groupStates, want) {
			g.logger.Info("FSM: want found in MainCity group", slog.String("state", want))
			return want, nil
		}
		g.logger.Warn("FSM: want not found in MainCity group, returning first element of group", slog.String("state", groupStates[0]))
		return groupStates[0], nil
	case vision.FuzzySubstringMatch(ocrFamily, "city", 1):
		g.logger.Debug("FSM: Family defined as WORLD (by city reference)",
			slog.String("ocr_family", ocrFamily),
			slog.String("want", want),
		)
		groupStates, ok := config.TitleToState["World"]
		g.logger.Debug("FSM: Group World", slog.Any("groupStates", groupStates), slog.Bool("found", ok))
		if ok && lo.Contains(groupStates, want) {
			g.logger.Info("FSM: want found in World group", slog.String("state", want))
			return want, nil
		}
		g.logger.Warn("FSM: want not found in World group, returning first element of group", slog.String("state", groupStates[0]))
		return groupStates[0], nil
	}

	// 1. Find state group by title (old path)
	groupStates, ok := getMatchedState(ocrTitle, 0)
	g.logger.Debug("FSM: getMatchedState result by ocrTitle",
		slog.String("ocr_title", ocrTitle),
		slog.Any("groupStates", groupStates),
		slog.Bool("found", ok),
	)
	if !ok {
		g.logger.Warn("FSM: Failed to match title with any state group",
			slog.String("ocr_title", ocrTitle),
			slog.String("expected_state", want),
		)
		return want, nil
	}

	// 2. Limit group by family (if possible)
	var filteredGroup []string
	switch {
	case vision.FuzzySubstringMatch(ocrFamily, "world", 1):
		filteredGroup = lo.Filter(groupStates, func(s string, _ int) bool {
			return s == state.StateMainCity
		})
		g.logger.Debug("FSM: Limiting by CITY family (MainCity)",
			slog.Any("filteredGroup", filteredGroup),
		)
	case vision.FuzzySubstringMatch(ocrFamily, "city", 1):
		filteredGroup = lo.Filter(groupStates, func(s string, _ int) bool {
			return s == state.StateWorld
		})
		g.logger.Debug("FSM: Limiting by WORLD family (World)",
			slog.Any("filteredGroup", filteredGroup),
		)
	default:
		filteredGroup = groupStates
		g.logger.Debug("FSM: Family not defined — using entire group",
			slog.Any("filteredGroup", filteredGroup),
		)
	}

	if lo.Contains(filteredGroup, want) {
		g.logger.Info("FSM: want found in filtered group", slog.String("state", want))
		return want, nil
	}

	if len(filteredGroup) > 0 {
		g.logger.Warn("FSM: want not found in filtered group, returning first element of filtered group",
			slog.String("state", filteredGroup[0]),
		)
		return filteredGroup[0], nil
	}

	g.logger.Warn("FSM: want not found in filtered group and filteredGroup is empty — returning want",
		slog.String("state", want),
	)

	return want, nil
}

func getMatchedState(title string, maxDistance int) ([]string, bool) {
	for key, states := range config.TitleToState {
		if vision.FuzzySubstringMatch(title, key, maxDistance) {
			return states, true
		}
	}
	return nil, false
}
