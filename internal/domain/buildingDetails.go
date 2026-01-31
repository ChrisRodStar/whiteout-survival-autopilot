package domain

import "time"

// BuildingDetails describes building parameters that are level-independent and rarely change.
// This data can be used to calculate construction time, cost, and building bonuses.
type BuildingDetails struct {
	// ConstructionTime – time required to construct the building (e.g., "2h30m").
	ConstructionTime time.Duration `yaml:"construction_time"`

	// Cost – resource costs for building construction.
	Cost Resources `yaml:"cost"`

	// Benefits – description of bonuses provided by the building (e.g., increased food production).
	Benefits string `yaml:"benefits"`
}
