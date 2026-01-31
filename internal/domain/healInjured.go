package domain

type HealInjured struct {
	State HealInjuredState `yaml:"state"` // Heal injured state
}

type HealInjuredState struct {
	IsAvailable bool   `yaml:"isAvailable"` // Flag indicating heal injured availability (by icon)
	IsNext      string `yaml:"isNext"`      // Flag indicating heal injured availability (by text)

	IsReplenishAll bool   `yaml:"isReplenishAll"` // Flag indicating availability to heal all injured
	StatusHeal     string `yaml:"statusHeal"`     // Heal injured status
}
