package settings

type Logic struct {
	Set               LogicSetting
	StartingTimeOfDay Enum
	BlueFireArrows    bool
	FixBrokenDrops    bool
	HintsRevealed     HintsRevealed
	FreeBombchuDrops  bool

	Spawns struct {
		StartTimeOfDay TimeOfDay
		AdultSpawn     string
		ChildSpawn     string
		StartAge       StartAge
	}

	Shuffling struct {
		DungeonRewards ShuffleDungeonRewards
		Songs          Enum
		Shops          Enum
		SkullTokens    ShuffleSkullTokens
		Scrubs         ShuffleScrub
		Freestandings  Enum
		Fishing        Flags
		Pots           ShufflePots
		Crates         ShuffleCrates
		Cows           bool

		RemoveCollectibleHearts bool

		Flags ShufflingFlags
	}

	Damage struct {
		Multiplier DamageMultiplier
		Bonks      DeadlyBonks
	}

	Locations struct {
		Reachability LocationsReachable
		Flags        Flags // rauru, zelda letter, mask quest, scarecrow, beans
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
		Interiors        InteriorShuffle
		DungeonDoor      DungeonDoorShuffle

		OpenDoorOfTime bool
		Grottos        bool
		Overworld      bool
		ValleyRiver    bool
		OwlDrops       bool
		WarpSongs      bool
	}

	WinConditions struct {
		TriforceHunt  bool
		TriforceCount int
		TriforceGoal  int
		Lacs          ConditionedAmount
		Bridge        ConditionedAmount
		GanonBossKey  ConditionedAmount
		Trials        Flags
	}

	Tricks map[string]bool
}

func finalizeLogic(_ *Zootr, _ *Logic) error {
	return notImpled
}
