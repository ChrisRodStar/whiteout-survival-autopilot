package fsm_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	fsm2 "github.com/batazor/whiteout-survival-autopilot/internal/fsm"
)

func TestExpectState_TableDriven(t *testing.T) {
	// Set ENV
	os.Setenv("PATH_TO_FSM_STATE_RULES", "../../references/fsmState.yaml")

	lookup, err := config.LoadAreaReferences("../../references/area.json")
	if err != nil {
		t.Fatalf("failed to load area.json: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	fakeADB := &FakeADB{}

	tests := []struct {
		name           string
		want           string
		ocrTitle       string
		ocrFamily      string
		expectedResult string
		expectedErr    bool
		analyzeErr     error
	}{
		{
			name:           "want in filtered group (city)",
			want:           "main_city",
			ocrTitle:       "MainCity", // <-- key matches the map!
			ocrFamily:      "world",    // filter on main_city
			expectedResult: "main_city",
		},
		{
			name:           "want in filtered group (world)",
			want:           "world",
			ocrTitle:       "MainCity",
			ocrFamily:      "city", // filter on world
			expectedResult: "world",
		},
		{
			name:           "want not in group, but group not empty — returns first",
			want:           "some_state",
			ocrTitle:       "MainCity",
			ocrFamily:      "",
			expectedResult: "some_state",
		},
		{
			name:           "no title matches — returns want",
			want:           "mail",
			ocrTitle:       "Unknown", // No such key
			ocrFamily:      "",
			expectedResult: "mail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTitleToState := config.TitleToState
			defer func() { config.TitleToState = oldTitleToState }()

			gamer := &domain.Gamer{
				ScreenState: domain.ScreenState{},
			}
			gameFSM := fsm2.NewGame(logger, fakeADB, lookup, nil, gamer)

			got, err := gameFSM.ExpectState(tt.want)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, got)
			}
		})
	}
}
