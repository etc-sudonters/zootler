package interpreter

//RuleParser.py
type Zoot_AtDay struct{}
type Zoot_AtDampeTime struct{}
type Zoot_AtNight struct{}

// State.py
// ("item name", qty) tuples and "raw_item_name" w/ implicit qty = 1, having more is fine
type Zoot_HasQuantityOf struct{}
type Zoot_HasAnyOf struct{}
type Zoot_HasAllOf struct{}
type Zoot_CountOf struct{}
type Zoot_HeartCount struct{}
type Zoot_HasHearts struct{}
type Zoot_HasMedallions struct{}
type Zoot_HasStones struct{}
type Zoot_HasDungeonRewards struct{}
type Zoot_HasOcarinaButtons struct{}
type Zoot_HasItemGoal struct{}
type Zoot_ItemCount struct{}
type Zoot_ItemNameCount struct{}
type Zoot_HasBottle struct{}
type Zoot_HasFullItemGoal struct{}
type Zoot_HasAllItemGoals struct{}
type Zoot_HadNightStart struct{}
type Zoot_CanLiveDmg struct{}
type Zoot_GuaranteeHint struct{}
type Zoot_RegionHasShortcuts struct{}
type Zoot_HasNotesForSong struct{}
