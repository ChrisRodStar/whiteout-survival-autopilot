package domain

import (
	"time"
)

type VIP struct {
	Level int           `yaml:"level"` // VIP status level (e.g., 1, 2, 3, etc.).
	Time  time.Duration `yaml:"time"`  // Time remaining until VIP status expires (e.g., 30 days).

	State VIPState `yaml:"state"` // VIP status state (e.g., active, expiring, etc.)
}

type VIPState struct {
	IsNotify           bool `yaml:"isNotify"`           // Flag indicating whether there are VIP status events.
	IsActive           bool `yaml:"isActive"`           // Flag indicating whether VIP status is active.
	IsAdd              bool `yaml:"isAdd"`              // Flag indicating whether VIP status can be added.
	IsAward            bool `yaml:"isAward"`            // Flag indicating whether VIP status reward is available.
	IsClaim            bool `yaml:"isClaim"`            // Flag indicating whether VIP status reward is available.
	IsVIPAddAvailable  bool `yaml:"isVIPAddAvailable"`  // Flag indicating whether the ability to add VIP status is available.
	IsVIPAddAvailableX bool `yaml:"isVIPAddAvailableX"` // Flag indicating whether the ability to add VIP status is available (additional flag).
}
