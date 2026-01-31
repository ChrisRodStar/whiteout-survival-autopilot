package domain

import (
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/heroes"
)

type Gamers []Gamer

// Gamer describes a game character with all characteristics.
type Gamer struct {
	ID       int    `yaml:"id"`       // Unique character identifier (fid).
	Nickname string `yaml:"nickname"` // Character nickname.
	State    int    `yaml:"state"`    // Character state.
	Avatar   string `yaml:"avatar"`   // Character avatar URL.
	Gems     int    `yaml:"gems"`     // Number of gems (premium currency).
	Power    int    `yaml:"power"`    // Character power.

	ScreenState ScreenState `yaml:"screenState"` // Screen state (e.g., "main", "battle", "exploration").

	VIP            VIP            `yaml:"vip"`            // Character VIP status.
	Resources      Resources      `yaml:"resources"`      // Character resources.
	Exploration    Exploration    `yaml:"exploration"`    // World exploration.
	Heroes         heroes.Heroes  `yaml:"heroes"`         // Heroes state.
	Messages       MessagesState  `yaml:"messages"`       // Messages state.
	Alliance       Alliance       `yaml:"alliance"`       // Alliance data.
	Buildings      Buildings      `yaml:"buildings"`      // Character buildings.
	Researches     Researches     `yaml:"researches"`     // Research levels.
	Events         Events         `yaml:"events"`         // Character events.
	Troops         Troops         `yaml:"troops"`         // Troops state.
	Tech           Tech           `yaml:"tech"`           // Character technologies.
	Mail           Mail           `yaml:"mail"`           // Character mail state.
	Shop           Shop           `yaml:"shop"`           // Shop state.
	DailyMissions  DailyMissions  `yaml:"dailyMissions"`  // Daily missions state.
	GrowthMissions GrowthMissions `yaml:"growthMissions"` // Character growth state.
	Chief          Chief          `yaml:"chief"`          // Governor data
	Arena          Arena          `yaml:"arena"`          // Arena data
	HealInjured    HealInjured    `yaml:"healInjured"`    // Heal injured events
}

// Len returns the number of gamers.
func (g Gamers) Len() int {
	return len(g)
}

// Swap exchanges the gamers at indices i and j.
func (g Gamers) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// Less compares two gamers by their Nickname.
// Adjust this comparison if you want to sort by another field (e.g., ID, Power, etc.).
func (g Gamers) Less(i, j int) bool {
	return g[i].Nickname < g[j].Nickname
}
