package domain

// Buildings represents a set of character buildings.
type Buildings struct {
	Queue1 string `yaml:"queue1"`
	Queue2 string `yaml:"queue2"`

	State BuildingState `yaml:"state"`

	Furnace Building `yaml:"furnace"` // Furnace.
}

type BuildingState struct {
	Text string `yaml:"text"`
}
