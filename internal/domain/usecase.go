package domain

import (
	"fmt"
	"time"
)

// UseCase represents a scenario described in a YAML file.
type UseCase struct {
	Name     string        `yaml:"name"`     // Scenario name
	Priority int           `yaml:"priority"` // Scenario priority (from 0 to 100)
	Node     string        `yaml:"node"`     // Initial screen/state from which the usecase starts
	Trigger  string        `yaml:"trigger"`  // CEL expression that determines whether to run the usecase
	Steps    Steps         `yaml:"steps"`    // Sequence of steps
	TTL      time.Duration `yaml:"ttl"`      // Usecase time-to-live (e.g., "24h")
	Cron     string        `yaml:"cron"`     // Cron expression for periodic usecase execution (e.g., "0 0 * * *")

	SourcePath string `json:"-"` // Path to the file from which the usecase was loaded
}

// Steps is simply a slice of Step
type Steps []Step

// Step represents an individual step in a usecase.
// It can be a simple action (click/wait), conditional (if), loop,
// include screenshot analysis (analyze), or manage TTL.
type Step struct {
	// Common actions
	Click   string        `yaml:"click,omitempty"`   // Name of the region to click
	Longtap string        `yaml:"longtap,omitempty"` // Name of the region to long tap
	Action  string        `yaml:"action,omitempty"`  // Special action: "loop", "screenshot", etc.
	Wait    time.Duration `yaml:"wait,omitempty"`    // Wait duration (e.g., "5s")

	// Conditional block if { then {} else {} }
	If *IfStep `yaml:"if,omitempty"`

	// Loop: used together with action: loop
	Trigger string `yaml:"trigger,omitempty"` // CEL expression used for loop or if
	Steps   Steps  `yaml:"steps,omitempty"`   // Nested steps (e.g., inside loop or if.then)

	// Screenshot analysis: used in pair with action: screenshot
	Analyze []AnalyzeRule `yaml:"analyze,omitempty"` // List of analysis rules (e.g., text/icon/etc.)

	// TTL management (e.g., to postpone repeated usecase execution)
	SetTTL      string `yaml:"setTTL,omitempty"`      // Duration, e.g., "24h"
	UsecaseName string `yaml:"usecaseName,omitempty"` // Name of the usecase to which TTL is applied

	Set string      `yaml:"set,omitempty"` // Path like "exploration.state.battleStatus"
	To  interface{} `yaml:"to,omitempty"`  // New value (in your case — "")

	PushUsecase []PushUsecase `yaml:"pushUsecase,omitempty"` // List of usecases to run when executing this step
}

// IfStep describes a conditional construct of the form if { then {} else {} }
type IfStep struct {
	Trigger string `yaml:"trigger"`        // CEL expression returning bool
	Then    []Step `yaml:"then"`           // Steps executed if trigger = true
	Else    []Step `yaml:"else,omitempty"` // Steps if trigger = false (optional)
}

// AnalyzeRule describes rules for analyzing a screen region (screenshot).
type AnalyzeRule struct {
	Name              string        `yaml:"name"`                        // Region name (and key for saving)
	Action            string        `yaml:"action"`                      // Action: "text", "exist", "color_check", "findIcon", "findText"
	Text              string        `yaml:"text,omitempty"`              // Text to search for (e.g., "Battle")
	Type              string        `yaml:"type,omitempty"`              // Result type (e.g., "integer" if action = text)
	Threshold         float64       `yaml:"threshold,omitempty"`         // Confidence level, default 0.9
	ExpectedColorBg   string        `yaml:"expectedColorBg,omitempty"`   // Expected background color (e.g., "red")
	ExpectedColorText string        `yaml:"expectedColorText,omitempty"` // Expected text color (e.g., "green")
	Log               string        `yaml:"log,omitempty"`               // Message for logging (optional)
	SaveAsRegion      bool          `yaml:"saveAsRegion,omitempty"`      // If true — save the zone as a new temporary region with name .Name
	PushUseCase       []PushUsecase `yaml:"pushUsecase,omitempty"`       // List of usecases to run when executing this rule
}

type PushUsecase struct {
	Trigger string    `yaml:"trigger"` // CEL expression
	List    []UseCase `yaml:"list"`    // Usecases to send to the queue
}

// Validate checks the validity of the action value in the analysis rule.
func (r AnalyzeRule) Validate() error {
	switch r.Action {
	case "text", "exist", "color_check", "findIcon", "findText":
		return nil
	default:
		return fmt.Errorf("invalid action '%s' in rule '%s'", r.Action, r.Name)
	}
}
