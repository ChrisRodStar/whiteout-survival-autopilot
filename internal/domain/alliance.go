package domain

// Alliance describes data about the alliance to which the character belongs.
type Alliance struct {
	Name    string        `yaml:"name"`    // Alliance name.
	MyLevel int           `yaml:"myLevel"` // R5-R1
	Power   int           `yaml:"power"`   // Alliance power.
	Members MembersInfo   `yaml:"members"` // Information about alliance members.
	State   AllianceState `yaml:"state"`   // Additional alliance state.

	// Left menu
	War       AllianceWar       `yaml:"war"`       // Alliance war (e.g., current war).
	Territory AllianceTerritory `yaml:"territory"` // Alliance territory (e.g., captured territories).
	Shop      AllianceShop      `yaml:"shop"`      // Alliance shop (e.g., purchases).

	// Right menu
	Chests AllianceChests `yaml:"chests"` // Alliance chests (e.g., rewards).
	Battle AllianceBattle `yaml:"battle"` // Alliance battle (e.g., battles with other alliances).
	Tech   AllianceTech   `yaml:"tech"`   // Alliance technologies (e.g., contribution).
	Help   AllianceHelp   `yaml:"help"`   // Alliance help (e.g., helping other players).
}

// MembersInfo contains information about the number of alliance members.
type MembersInfo struct {
	Count int `yaml:"count"` // Current number of members.
	Max   int `yaml:"max"`   // Maximum number of members.
}

// AllianceState describes the alliance state.
type AllianceState struct {
	IsNeedSupport              bool `yaml:"isNeedSupport"`              // Flag for alliance participation.
	IsWar                      int  `yaml:"isWar"`                      // Number of current wars.
	IsChests                   int  `yaml:"isChests"`                   // Number of available chests.
	IsAllianceContributeButton bool `yaml:"isAllianceContributeButton"` // Technology contribution button
	IsAllianceTechButton       bool `yaml:"isAllianceTechButton"`       // Alliance technologies button
	PolarTerrorCount           int  `yaml:"polarTerrorCount"`           // Number of successful polar bear joins

	// chests
	IsClaimButton        bool `yaml:"isClaimButton"`        // Button to claim alliance reward
	IsCanClaimAllChests  bool `yaml:"isCanClaimAllChests"`  // Button to claim all chests
	LootCountLimit       int  `yaml:"lootCountLimit"`       // Chest limit
	IsGiftClaimAllButton bool `yaml:"isGiftClaimAllButton"` // Button to claim all gifts
	IsMainChest          bool `yaml:"isMainChest"`          // Button to claim main chest
}

// AllianceWar describes the alliance war. --------------------------------------
type AllianceWar struct {
	IsNotify bool `yaml:"isNotify"` // Flag for war notification.
}

// AllianceTerritory describes the alliance territory.
type AllianceTerritory struct {
	IsNotify bool `yaml:"isNotify"` // Flag for territory notification.
}

// AllianceShop describes the alliance shop.
type AllianceShop struct{}

// AllianceChests describes the alliance chests. -----------------------------
type AllianceChests struct {
	IsNotify bool `yaml:"isNotify"` // Flag for chest notification.
}

// AllianceBattle describes the alliance battles.
type AllianceBattle struct {
	IsNotify bool `yaml:"isNotify"` // Flag for battle notification.
}

// AllianceTech describes the technological aspects of the alliance.
type AllianceTech struct {
	IsNotify bool `yaml:"isNotify"` // Flag for technology notification.

	Favorite bool `yaml:"favorite"` // Flag for technology contribution.
}

// AllianceHelp describes the alliance help.
type AllianceHelp struct {
	IsNotify bool `yaml:"isNotify"` // Flag for help notification.
}
