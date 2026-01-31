package domain

// State contains the list of accounts that the bot works with.
type State struct {
	Gamers Gamers `yaml:"gamers"` // Game characters
}

// Resources describes the resources owned by a character.
type Resources struct {
	Wood int `yaml:"wood"` // Wood.
	Food int `yaml:"food"` // Food.
	Iron int `yaml:"iron"` // Iron.
	Meat int `yaml:"meat"` // Meat.
}

// MessagesState contains information about character messages.
type MessagesState struct {
	State MessageStatus `yaml:"state"`
}

// MessageStatus describes the message status.
type MessageStatus struct {
	IsNewMessage bool `yaml:"isNewMessage"` // Flag indicating the presence of new messages.
	IsNewReports bool `yaml:"isNewReports"` // Flag indicating the presence of new reports.
}

// Building describes an individual building.
type Building struct {
	Level int `yaml:"level"` // Building level.
	Power int `yaml:"power"` // Building power.
	// Additional fields can be added, such as construction time, resource costs, etc.
}

// Researches describes the character's research levels.
type Researches struct {
	Battle  Research `yaml:"battle"`  // Military research.
	Economy Research `yaml:"economy"` // Economic research.
	// Additional research can be added here.
}

// Research describes the level of a specific research.
type Research struct {
	Level int `yaml:"level"` // Research level.
	// Additional fields if needed.
}
