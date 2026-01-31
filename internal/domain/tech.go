package domain

type Tech struct {
	State TechState `yaml:"state"` // Technology state.
}

type TechState struct {
	IsAvailable bool   `yaml:"is_available"` // Technology availability.
	TextStatus  string `yaml:"TextStatus"`   // Text status.
}
