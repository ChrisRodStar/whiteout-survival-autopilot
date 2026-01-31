package state

// --------------------------------------------------------------------
// State Definitions: Each constant represents a game screen (state)
// --------------------------------------------------------------------
const (
	InitialState         = "initial"
	StateMainCity        = "main_city"
	StateActivityTriumph = "activity_triumph"
	StateProfile         = "profile"
	StateLeaderboard     = "leaderboard"
	StateSettings        = "settings"
	StateDawnMarket      = "dawn_market"

	// Pets
	StatePets = "pets"

	// Exploration
	StateExploration       = "exploration"
	StateExplorationBattle = "exploration_battle"

	// Account switching
	StateChiefProfile                           = "chief_profile"
	StateChiefCharacters                        = "chief_characters"
	StateChiefProfileSetting                    = "chief_profile_setting"
	StateChiefProfileAccount                    = "chief_profile_account"
	StateChiefProfileAccountChangeAccount       = "chief_profile_account_change_account"
	StateChiefProfileAccountChangeGoogle        = "chief_profile_account_change_account_google"
	StateChiefProfileAccountChangeGoogleConfirm = "chief_profile_account_change_account_google_continue"

	StateAllianceWar              = "alliance_war"
	StateAllianceWarRally         = "alliance_war_rally"
	StateAllianceWarRallyAutoJoin = "alliance_war_rally_auto_join"
	StateAllianceWarSolo          = "alliance_war_solo"
	StateAllianceWarEvents        = "alliance_war_events"

	// Global map
	StateWorld          = "world"
	StateWorldSearch    = "world_search_resources"
	StateWorldGlobalMap = "world_global_map"

	// Messages
	StateMail         = "mail"
	StateMailWars     = "mail_wars"
	StateMailAlliance = "mail_alliance"
	StateMailSystem   = "mail_system"
	StateMailReports  = "mail_reports"
	StateMailStarred  = "mail_starred"

	// VIP
	StateVIP    = "vip"
	StateVIPAdd = "vip_add"

	// Governor
	StateChiefOrders = "chief_orders"

	// Daily missions
	StateDailyMissions = "daily_missions"
	// Growth missions
	StateGrowthMissions = "growth_missions"
)

const (
	// Alliance
	StateAllianceManage    = "alliance_manage"
	StateAllianceTech      = "alliance_tech"
	StateAllianceSettings  = "alliance_settings"
	StateAllianceRanking   = "alliance_ranking"
	StateAllianceTerritory = "alliance_territory"

	// Alliance - chests
	StateAllianceChests    = "alliance_chests"
	StateAllianceChestLoot = "alliance_chest_loot"
	StateAllianceChestGift = "alliance_chest_gift"

	// Triumph
	StateAllianceActivityTriumph = "alliance_activity_triumph"
)

const (
	// Tundra Adventure
	StateTundraAdventure               = "tundra_adventure"
	StateTundraAdventureMain           = "tundra_adventure_main"
	StateTundraAdventureDrill          = "tundra_adventure_drill"
	StateTundraAdventurerDrill         = "tundra_adventurer_drill"
	StateTundraAdventurerDailyMissions = "tundra_adventurer_daily_missions"
	StateTundraAdventureOdessey        = "tundra_adventure_odessey"
	StateTundraAdventureCaravan        = "tundra_adventure_caravan"
)

const (
	StateInfantryCityView = "infantry_city_view"
	StateLancerCityView   = "lancer_city_view"
	StateMarksmanCityView = "marksman_city_view"
)

const (
	// Main menu
	StateMainMenuCity         = "main_menu_city"
	StateMainMenuWilderness   = "main_menu_wilderness"
	StateMainMenuBuilding1    = "main_menu_building_1"
	StateMainMenuBuilding2    = "main_menu_building_2"
	StateMainMenuTechResearch = "main_menu_tech_research"
)

const (
	// Backpack
	StateBackpack          = "backpack"
	StateBackpackResources = "backpack_resources"
	StateBackpackSpeedups  = "backpack_speedups"
	StateBackpackBonus     = "backpack_bonus"
	StateBackpackGear      = "backpack_gear"
	StateBackpackOther     = "backpack_other"
)

const (
	// Chat
	StateChat         = "chat"
	StateChatAlliance = "chat_alliance"
	StateChatWorld    = "chat_world"
	StateChatPersonal = "chat_personal"
)

const (
	// Heroes
	StateHeroes = "heroes"
)

const (
	// Events
	StateEvents = "events"
)

const (
	StateDeals       = "deals"
	StateTopUpCenter = "top_up_center"
)

const (
	// Intelligence
	StateIntel = "intel"
)

const (
	StateArenaCityView             = "arena_city_view"
	StateArenaMain                 = "arena_main"
	StateArenaDefensiveSquadLineup = "arena_defensive_squad_lineup"
	StateArenaChallengeList        = "arena_challenge_list"
)

const (
	StateFishingCityView = "fishing_city_view"
	StateFishingMain     = "fishing_main"
)

const (
	// Healing
	StateHealInjured = "heal_injured"
)

const (
	// Labyrinth
	StateLabyrinth      = "labyrinth"
	StateCaveOfMonsters = "cave_of_monsters"
)

const (
	// Heroes
	StateNatalia = "heroes_natalia"
)

const (
	// Enlistment Office
	StateEnlistmentOffice = "enlistment_office"
)
