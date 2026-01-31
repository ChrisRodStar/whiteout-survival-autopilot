package domain

type DailyMissions struct {
	IsNotify bool `yaml:"isNotify"` // Flag indicating whether daily mission rewards are available.

	State DailyMissionsState `yaml:"state"` // Daily missions state.

	Tasks Tasks `yaml:"tasks"` // Tasks that need to be completed.
}

type Tasks struct {
	IsReseachOneTechnologies bool `yaml:"isReseachOneTechnologies"` // Flag indicating whether one technology research is completed.
	IsGatherMeat             bool `yaml:"isGatherMeat"`             // Flag indicating whether the meat gathering task is completed.
}

type DailyMissionsState struct {
	IsClaimAll    bool `yaml:"isClaimAll"`    // Flag indicating whether the task to claim all rewards is completed.
	IsClaimButton bool `yaml:"isClaimButton"` // Flag indicating whether the claim reward button is available.
}

type GrowthMissions struct {
	IsNotify bool `yaml:"isNotify"` // Flag indicating whether growth mission rewards are available.

	State GrowthMissionsState `yaml:"state"` // Growth missions state.
}

type GrowthMissionsState struct {
	IsClaimAll    bool `yaml:"isClaimAll"`    // Flag indicating whether the task to claim all rewards is completed.
	IsClaimButton bool `yaml:"isClaimButton"` // Flag indicating whether the claim reward button is available.
}
