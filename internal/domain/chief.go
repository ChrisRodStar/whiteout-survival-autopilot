package domain

type Chief struct {
	Contentment int `yaml:"contentment"` // Governor satisfaction points

	State ChiefState `yaml:"state"` // Governor state
}

type ChiefState struct {
	IsNotify bool `yaml:"isNotify"` // Flag indicating whether there are governor orders.

	IsUrgentMobilization bool `yaml:"isUrgentMobilization"` // Flag indicating whether there is urgent mobilization.
	IsComprehensiveCare  bool `yaml:"isComprehensiveCare"`  // Flag indicating whether there is comprehensive care.
	IsProductivityDay    bool `yaml:"isProductivityDay"`    // Flag indicating whether there is productivity day.
	IsRushJob            bool `yaml:"isRushJob"`            // Flag indicating whether there is rush job.
	IsDoubleTime         bool `yaml:"isDoubleTime"`         // Flag indicating whether there is double time.
	IsFestivities        bool `yaml:"isFestivities"`        // Flag indicating whether there are festivities.
}
