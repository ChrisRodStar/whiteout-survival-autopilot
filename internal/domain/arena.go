package domain

type Arena struct {
	Rank    int `yaml:"rank"`    // Player's rank in the arena
	MyPower int `yaml:"myPower"` // Player's power

	State ArenaState `yaml:"state"` // Arena state (e.g., "open", "closed", "in_battle").
}

type ArenaState struct {
	IsFreeRefresh       bool `yaml:"isFreeRefresh"`       // Flag for free opponent refresh availability.
	IsAvailableFight    bool `yaml:"isAvailableFight"`    // Flag for fight availability.
	CountAvailableFight int  `yaml:"countAvailableFight"` // Number of available fights.

	EnemyPower1 int `yaml:"enemyPower1"` // First opponent's power.
	EnemyPower2 int `yaml:"enemyPower2"` // Second opponent's power.
	EnemyPower3 int `yaml:"enemyPower3"` // Third opponent's power.
	EnemyPower4 int `yaml:"enemyPower4"` // Fourth opponent's power.
	EnemyPower5 int `yaml:"enemyPower5"` // Fifth opponent's power.
}
