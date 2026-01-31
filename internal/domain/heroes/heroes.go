package heroes

// Heroes contains information about the state of all heroes.
type Heroes struct {
	IsNotify bool `json:"isNotify"` // Flag indicating the need to notify about hero state.

	List map[string]Hero
}

// Hero represents a single hero with their characteristics and state.
type Hero struct {
	Class      string            `json:"class"`           // Infantry, Lancer, Marksman
	Generation int               `json:"generation"`      // Hero generation
	Roles      []string          `json:"roles"`           // Roles (rally_leader, resource_gathering, etc.)
	Skills     HeroSkills        `json:"skills"`          // Hero skills
	Buffs      map[string]string `json:"buffs,omitempty"` // Buffs, key-value
	Notes      string            `json:"notes,omitempty"` // Notes on hero usage

	State State `json:"state,omitempty"` // Current hero state for the user
}

// HeroSkills contains skill groups (currently — only expedition).
type HeroSkills struct {
	Expedition map[string]HeroSkill `json:"expedition"` // one, two, three, etc.
}

// HeroSkill represents a specific skill.
type HeroSkill struct {
	Name     string `json:"name"`     // Skill name
	Priority int    `json:"priority"` // Skill priority (higher — more important)
}

// State describes the current hero status for the user.
type State struct {
	Level         int  `json:"level"`           // Upgrade level
	IsAvailable   bool `json:"is_available"`    // Whether the hero is available to the user
	IsCampTrainer bool `json:"is_camp_trainer"` // Whether the hero is a camp trainer
}
