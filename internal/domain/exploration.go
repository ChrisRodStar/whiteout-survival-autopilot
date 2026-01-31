package domain

// Exploration describes the world exploration level.
type Exploration struct {
	Level int              `yaml:"level"`
	State ExplorationState `yaml:"state"`

	IsNotify bool `yaml:"isNotify"` // Flag indicating the need to notify about exploration state.
}

// ExplorationState describes the world exploration state.
type ExplorationState struct {
	IsClaimActive bool `yaml:"isClaimActive"` // Flag indicating "Claim" button availability.

	MyPower      int    `yaml:"myPower"`      // Player's power.
	EnemyPower   int    `yaml:"enemyPower"`   // Enemy's power.
	BattleStatus string `yaml:"battleStatus"` // Battle status (e.g., "victory", "defeat").
}
