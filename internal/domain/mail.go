package domain

type Mail struct {
	IsHasMail int `yaml:"isHasMail"` // Flag indicating that mail exists.

	State MailState `yaml:"state"` // Mail state.
}

type MailState struct {
	IsWars     int `yaml:"isWars"`     // Flag indicating that mail contains war information.
	IsAlliance int `yaml:"isAlliance"` // Flag indicating that mail contains alliance information.
	IsSystem   int `yaml:"isSystem"`   // Flag indicating that mail contains system messages.
	IsReports  int `yaml:"isReports"`  // Flag indicating that mail contains reports.
}
