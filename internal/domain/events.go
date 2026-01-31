package domain

type Events struct {
	TundraAdventure TundraAdventure `yaml:"tundraAdventure"` // Tundra events
	FrostyFortune   FrostyFortune   `yaml:"frostyFortune"`   // Frosty fortress events
}

type TundraAdventure struct {
	State TundraAdventureState `yaml:"state"` // Tundra state
}

type TundraAdventureState struct {
	// Main City ------------
	IsExist bool `yaml:"isExist"` // Flag indicating event existence

	// Play screen -----------
	Count  int  `yaml:"count"`  // Number of available rolls
	IsPlay bool `yaml:"isPlay"` // Flag indicating dice roll availability

	// Adventurer Drill ------
	IsAdventurerDrillClaimIsExist bool `yaml:"isAdventurerDrillClaimIsExist"` // Flag indicating data update existence
	IsAdventurerDrillClaim        bool `yaml:"isAdventurerDrillClaim"`        // Flag indicating quest to claim loot

	// Daily Missions Rewards ------
	IsAdventureDailyClaim bool `yaml:"isAdventureDailyClaim"` // Flag indicating data update existence
}

type FrostyFortune struct {
	State FrostyFortuneState `yaml:"state"` // Frosty fortress state
}

type FrostyFortuneState struct {
	IsExist bool `yaml:"isExist"` // Flag indicating event existence
}
