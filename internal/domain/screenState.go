package domain

import (
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
)

type ScreenState struct {
	IsMainMenu bool   `yaml:"isMainMenu"` // Flag indicating whether there are events in the main menu.
	IsWelcome  bool   `yaml:"isWelcome"`  // Flag indicating whether there are welcome events for new survivors.
	IsMainCity string `yaml:"isMainCity"` // Flag indicating which screen is active - city or world map.

	CurrentState string `yaml:"currentState"` // Screen title.
	TitleFact    string `yaml:"titleFact"`    // Screen title obtained from screenshot analysis.
}

// Reset resets the screen state.
func (s *ScreenState) Reset() {
	s.IsMainMenu = false
	s.IsWelcome = false
	s.IsMainCity = ""
	s.CurrentState = state.StateMainCity
	s.TitleFact = ""
}
