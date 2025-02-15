package settings

type Logic struct {
	Set               LogicSetting
	StartingTimeOfDay Enum
	BlueFireArrows    bool
	FixBrokenDrops    bool
	HintsRevealed     HintsRevealed
	FreeBombchuDrops  bool
	StartingHearts    int

	Spawns struct {
		StartTimeOfDay TimeOfDay
		AdultSpawn     string
		ChildSpawn     string
		StartAge       StartAge
	}

	Shuffling struct {
		DungeonRewards ShuffleDungeonRewards
		Songs          ShuffleSongs
		Shops          ShuffleShop
		SkullTokens    ShuffleSkullTokens
		Scrubs         ShuffleScrub
		Freestandings  ShuffleFreestanding
		Pots           ShufflePots
		Crates         ShuffleCrates

		RemoveCollectibleHearts bool

		Flags ShufflingFlags
	}

	Damage struct {
		Multiplier DamageMultiplier
		Bonks      DeadlyBonks
	}

	Locations struct {
		Reachability LocationsReachable
		Flags        LocationFlags // rauru, zelda letter, mask quest, scarecrow, beans
		AdultSpawn   string
		ChildSpawn   string

		Disabled []string
	}

	Trade struct {
		AdultItems AdultTradeItems
		ChildItems ChildTradeItems
	}

	Dungeon struct {
		OneMajorItemPerDungeon bool
		Bosses                 BossShuffle
		Keys                   ShuffleKeys
		BossKey                ShuffleKeys
		SilverRupees           ShuffleKeys
		KeyRings               bool
		KeyRingsIncludeBossKey bool
		SilverRupeePouches     bool
		MapCompass             NavigationShuffle
		Empty                  Flags
		MasterQuest            Flags
		Shortcuts              Flags
		GerudoFortressKeys     ShuffleKeys
		GerudoFortress         GerudoFortressCarpenterRescue
		GanonBossKeyShuffle    GanonBossKeyShuffle
	}

	Minigames struct {
		TreasureChestGameKeys ShuffleKeys
		KakarikoChickenGoal   uint8
		BigPoeGoal            uint8
	}

	Connections struct {
		OpenKakarikoGate OpenKakarikoGate
		OpenKokriForest  OpenForest
		OpenZoraFountain OpenZoraFountain

		Flags     ConnectionFlag
		Interior  InteriorShuffle
		Overworld OverworldShuffle
	}

	WinConditions struct {
		TriforceHunt  bool
		TriforceCount int
		TriforceGoal  int
		Lacs          ConditionedAmount
		Bridge        ConditionedAmount
		GanonBossKey  ConditionedAmount
		Trials        TrialFlag
	}

	Tricks map[string]bool
}

func finalizeLogic(_ *Zootr, _ *Logic) error {
	return notImpled
}
