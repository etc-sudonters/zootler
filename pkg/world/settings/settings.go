package settings

type SeedSettings struct {
	Logic    LogicRuleSet
	ItemPool ItemPool
	LogicSettings
	ShuffleSettings
}

type LogicSettings struct {
	KokriForest     KokiriForest
	KakGate         KakarikoGate
	DoorOfTime      DoorOfTime
	Fountain        ZorasFountain
	Bridge          BridgeRequirement
	TowerTrials     TowerTrialCount
	StartingAge     StartingAge
	ChildTradeQuest ChildTradeQuest
	AdultTradeItems AdultTradeItems
}

type ShuffleSettings struct {
	ShuffleSongs            SongShuffle
	ShuffleShops            ShopShuffle
	ShuffleTokens           GoldTokenShuffle
	ShuffleScrubs           ScrubShuffle
	ShufflePots             PotShuffle
	ShuffleCrate            CrateShuffle
	ShuffleCows             CowShuffle
	ShuffleBeehinves        BeehiveShuffle
	ShuffleKokriSword       KokriSwordShuffle
	ShuffleOcarinas         OcarinaShuffle
	ShuffleGerudoCard       GerudoCardShuffle
	ShuffleMagicBeans       MagicBeanShuffle
	ShuffleRepeatMerchants  RepeatMerchantShuffle
	ShuffleFrogRupees       FrogRupeeShuffle
	ShuffleMapsAndCompasses MapsAndCompassesShuffle
	ShuffleSmallKeys        SmallKeyShuffle
	ShuffleBossKeys         BossKeyShuffle
	ShuffleTowerBossKey     TowerBossKeyShuffle
	ShuffleChestGameKeys    ChestGameKeyShuffle
}
