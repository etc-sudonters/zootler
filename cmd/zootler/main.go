package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"
)

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

type cliOptions struct {
	logicDir string
	dataDir  string
}

func (opts *cliOptions) init() {
	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.Parse()
}

func (c cliOptions) validate() error {
	if c.logicDir == "" {
		return missingRequired("-l")
	}

	if c.dataDir == "" {
		return missingRequired("-d")
	}

	return nil
}

type std struct{ *dontio.Std }

func (s std) WriteLineOut(msg string, v ...any) {
	fmt.Fprintf(s.Out, msg+"\n", v...)
}

func main() {
	var opts cliOptions
	var appExitCode stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	defer func() {
		os.Exit(int(appExitCode))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(stdio.Err, "%s\n", err)
			}
			_, _ = stdio.Err.Write(debug.Stack())
			if appExitCode != stageleft.ExitSuccess {
				appExitCode = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	exitWithErr := func(code stageleft.ExitCode, err error) {
		appExitCode = code
		fmt.Fprintf(stdio.Err, "%s\n", err.Error())
		panic("aaaH!")
	}

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	(&opts).init()

	if cliErr := opts.validate(); cliErr != nil {
		exitWithErr(2, cliErr)
		return
	}

	app, err := app.NewApp(ctx,
		app.ConfigureStorage(CreateScheme{DDL: []DDL{
			BitColumnOf[components.Advancement],
			BitColumnOf[components.Beehive],
			BitColumnOf[components.BossHeart],
			BitColumnOf[components.BossKey],
			BitColumnOf[components.Boss],
			BitColumnOf[components.Bottle],
			BitColumnOf[components.BottomoftheWellMQ],
			BitColumnOf[components.BottomoftheWell],
			BitColumnOf[components.Chest],
			BitColumnOf[components.CollectableGameToken],
			BitColumnOf[components.Collectable],
			BitColumnOf[components.Collected],
			BitColumnOf[components.Compass],
			BitColumnOf[components.Cows],
			BitColumnOf[components.Crate],
			BitColumnOf[components.Cutscene],
			BitColumnOf[components.DeathMountainCrater],
			BitColumnOf[components.DeathMountainTrail],
			BitColumnOf[components.DeathMountain],
			BitColumnOf[components.DekuScrubUpgrades],
			BitColumnOf[components.DekuScrubs],
			BitColumnOf[components.DekuTreeMQ],
			BitColumnOf[components.DekuTree],
			BitColumnOf[components.DesertColossus],
			BitColumnOf[components.DodongosCavernMQ],
			BitColumnOf[components.DodongosCavern],
			BitColumnOf[components.Drop],
			BitColumnOf[components.DungeonReward],
			BitColumnOf[components.Event],
			BitColumnOf[components.FireTempleMQ],
			BitColumnOf[components.FireTemple],
			BitColumnOf[components.FlyingPot],
			BitColumnOf[components.ForestArea],
			BitColumnOf[components.ForestTempleMQ],
			BitColumnOf[components.ForestTemple],
			BitColumnOf[components.Forest],
			BitColumnOf[components.Freestanding],
			BitColumnOf[components.GanonBossKey],
			BitColumnOf[components.GanonsCastleMQ],
			BitColumnOf[components.GanonsCastle],
			BitColumnOf[components.GanonsTower],
			BitColumnOf[components.GerudoTrainingGroundMQ],
			BitColumnOf[components.GerudoTrainingGround],
			BitColumnOf[components.GerudoValley],
			BitColumnOf[components.Gerudo],
			BitColumnOf[components.GerudosFortress],
			BitColumnOf[components.GoldSkulltulaToken],
			BitColumnOf[components.GoldSkulltulas],
			BitColumnOf[components.GoronCity],
			BitColumnOf[components.Graveyard],
			BitColumnOf[components.GreatFairies],
			BitColumnOf[components.GrottoScrub],
			BitColumnOf[components.Grottos],
			BitColumnOf[components.HauntedWasteland],
			BitColumnOf[components.HideoutSmallKey],
			BitColumnOf[components.HintStone],
			BitColumnOf[components.Hint],
			BitColumnOf[components.HyruleCastle],
			BitColumnOf[components.HyruleField],
			BitColumnOf[components.IceCavernMQ],
			BitColumnOf[components.IceCavern],
			BitColumnOf[components.Item],
			BitColumnOf[components.JabuJabusBellyMQ],
			BitColumnOf[components.JabuJabusBelly],
			BitColumnOf[components.Junk],
			BitColumnOf[components.KakarikoVillage],
			BitColumnOf[components.Kakariko],
			BitColumnOf[components.KokiriForest],
			BitColumnOf[components.LakeHylia],
			BitColumnOf[components.Location],
			BitColumnOf[components.LonLonRanch],
			BitColumnOf[components.LostWoods],
			BitColumnOf[components.Map],
			BitColumnOf[components.Market],
			BitColumnOf[components.MaskShop],
			BitColumnOf[components.MasterQuest],
			BitColumnOf[components.Medallion],
			BitColumnOf[components.Minigames],
			BitColumnOf[components.NPC],
			BitColumnOf[components.NeedSpiritualStones],
			BitColumnOf[components.OutsideGanonsCastle],
			BitColumnOf[components.Pot],
			BitColumnOf[components.Priority],
			BitColumnOf[components.Refill],
			BitColumnOf[components.RupeeTower],
			BitColumnOf[components.SacredForestMeadow],
			BitColumnOf[components.Scrub],
			BitColumnOf[components.ShadowTempleMQ],
			BitColumnOf[components.ShadowTemple],
			BitColumnOf[components.Shop],
			BitColumnOf[components.SilverRupee],
			BitColumnOf[components.SkulltulaHouse],
			BitColumnOf[components.SmallCrate],
			BitColumnOf[components.SmallKey],
			BitColumnOf[components.SpiritTempleMQ],
			BitColumnOf[components.SpiritTemple],
			BitColumnOf[components.SpiritualStone],
			BitColumnOf[components.TCGSmallKey],
			BitColumnOf[components.TempleofTime],
			BitColumnOf[components.ThievesHideout],
			BitColumnOf[components.Trade],
			BitColumnOf[components.VanillaDungeons],
			BitColumnOf[components.WaterTempleMQ],
			BitColumnOf[components.WaterTemple],
			BitColumnOf[components.Wonderitem],
			BitColumnOf[components.ZorasDomain],
			BitColumnOf[components.ZorasFountain],
			BitColumnOf[components.ZorasRiver],
			MapColumn[components.Count],
			MapColumn[components.Price],
			MapColumn[components.ShopObject],
			MapColumn[components.OcarinaButton],
			MapColumn[components.OcarinaNote],
			MapColumn[components.OcarinaSong],
			MapColumn[components.Song],
			SliceColumn[components.DefaultItem],
			SliceColumn[components.Inhabited],
			SliceColumn[components.Inhabits],
			SliceColumn[components.Name],
		}}),
		app.ConfigureStorage(DataFileLoader[FileItem]("inputs/data/items.json")),
		app.ConfigureStorage(DataFileLoader[FileLocation]("inputs/data/locations.json")),
	)

	if err != nil {
		exitWithErr(3, err)
		return
	}

	if err := example(app.Ctx(), app.Engine()); err != nil {
		exitWithErr(4, err)
	}
}
