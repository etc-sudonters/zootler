package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/stageleft"
	"muzzammil.xyz/jsonc"
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

func main() {
	var opts cliOptions
	var exit stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	defer func() {
		os.Exit(int(exit))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(stdio.Err, "%s\n", err)
			}
			_, _ = stdio.Err.Write(debug.Stack())
			if exit != stageleft.ExitSuccess {
				exit = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	(&opts).init()

	if cliErr := opts.validate(); cliErr != nil {
		fmt.Fprintf(stdio.Err, "%s\n", cliErr.Error())
		exit = stageleft.ExitCode(2)
		return
	}

	storage := query.NewEngine()
	buildStorage(ctx, storage)
	loadLocations("inputs/data/locations.json", storage)
	loadItems("inputs/data/items.json", storage)
	example(ctx, storage)

}

func example(ctx context.Context, storage query.Engine) {
	stdio, _ := dontio.StdFromContext(ctx)
	q := storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	allLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all locations: %d\n", allLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song locations: %d\n", songLocs.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Location]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongLocs, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of not Song locations: %d\n", notSongLocs.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	allToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of all tokens: %d\n", allToks.Len())

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	q.Exists(mirrors.TypeOf[components.Song]())
	songToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of Song tokens: %d\n", songToks.Len())

	q = storage.CreateQuery()

	q = storage.CreateQuery()
	q.Exists(mirrors.TypeOf[components.Token]())
	q.NotExists(mirrors.TypeOf[components.Song]())
	notSongToks, err := storage.Retrieve(q)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(stdio.Out, "Count of not Song tokens: %d\n", notSongToks.Len())
}

func loadLocations(path string, storage query.Engine) {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var locs []struct {
		Name       string   `json:"name"`
		Type       string   `json:"type"`
		Default    string   `json:"vanilla"`
		Categories []string `json:"categories"`
	}

	if err := jsonc.Unmarshal(raw, &locs); err != nil {
		panic(err)
	}

	var song components.Song
	var location components.Location

	for _, l := range locs {
		id, err := storage.InsertRow(components.Name(l.Name), location)
		if err != nil {
			panic(err)
		}
		if l.Type == "Song" {
			storage.SetValues(id, table.Values{song})
		}
	}
}

func loadItems(path string, storage query.Engine) {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var items []struct {
		Name        string                 `json:"name"`
		Type        string                 `json:"type"`
		Advancement bool                   `json:"advancement"`
		Priority    bool                   `json:"priority"`
		Special     map[string]interface{} `json:"special"`
	}

	var tok components.Token
	var song components.Song

	if err := jsonc.Unmarshal(raw, &items); err != nil {
		panic(err)
	}

	for _, item := range items {
		id, err := storage.InsertRow(components.Name(item.Name), tok)
		if err != nil {
			panic(err)
		}
		if item.Type == "Song" {
			storage.SetValues(id, table.Values{song})
		}
	}

}

func buildStorage(ctx context.Context, storage query.Engine) error {
	storage.CreateColumn(table.BuildColumnOf[components.Name](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.DefaultItem](columns.NewSlice()))

	storage.CreateColumn(table.BuildColumnOf[components.Collectable](columns.NewBit(components.Collectable{})))
	storage.CreateColumn(table.BuildColumnOf[components.Collected](columns.NewBit(components.Collected{})))
	storage.CreateColumn(table.BuildColumnOf[components.Inhabited](columns.NewSlice()))
	storage.CreateColumn(table.BuildColumnOf[components.Inhabits](columns.NewSlice()))

	storage.CreateColumn(table.BuildColumnOf[components.Alias](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.BossKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Bottle](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.BottomoftheWell](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Compass](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Count](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DeathMountainCrater](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DeathMountainTrail](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuScrubUpgrades](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuScrubs](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DekuTree](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DesertColossus](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DodongosCavern](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.DungeonReward](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.FireTemple](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ForestArea](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ForestTemple](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GanonBossKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GanonsCastle](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GanonsTower](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GerudoTrainingGround](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GerudoValley](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GerudosFortress](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GoldSkulltulaToken](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GoldSkulltulas](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GoronCity](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Graveyard](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.GreatFairies](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Grottos](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.HauntedWasteland](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.HideoutSmallKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.HyruleCastle](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.HyruleField](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.IceCavern](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.JabuJabusBelly](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.KakarikoVillage](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.KokiriForest](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.LakeHylia](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.LonLonRanch](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.LostWoods](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Map](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Market](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Medallion](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Minigames](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.NPCs](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.NeedSpiritualStones](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.OcarinaButton](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.OutsideGanonsCastle](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Price](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SacredForestMeadow](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ShadowTemple](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ShopObject](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Shops](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SkulltulaHouse](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SmallKey](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Song](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Spawn](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SpiritTemple](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.SpiritualStone](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.TempleofTime](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ThievesHideout](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.Trade](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.WaterTemple](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasDomain](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasFountain](columns.NewMap()))
	storage.CreateColumn(table.BuildColumnOf[components.ZorasRiver](columns.NewMap()))

	storage.CreateColumn(table.BuildColumnOf[components.Beehives](columns.NewBit(components.Beehives{})))
	storage.CreateColumn(table.BuildColumnOf[components.Chests](columns.NewBit(components.Chests{})))
	storage.CreateColumn(table.BuildColumnOf[components.Cows](columns.NewBit(components.Cows{})))
	storage.CreateColumn(table.BuildColumnOf[components.Crates](columns.NewBit(components.Crates{})))
	storage.CreateColumn(table.BuildColumnOf[components.Drop](columns.NewBit(components.Drop{})))
	storage.CreateColumn(table.BuildColumnOf[components.Event](columns.NewBit(components.Event{})))
	storage.CreateColumn(table.BuildColumnOf[components.FlyingPots](columns.NewBit(components.FlyingPots{})))
	storage.CreateColumn(table.BuildColumnOf[components.Freestandings](columns.NewBit(components.Freestandings{})))
	storage.CreateColumn(table.BuildColumnOf[components.Item](columns.NewBit(components.Item{})))
	storage.CreateColumn(table.BuildColumnOf[components.Junk](columns.NewBit(components.Junk{})))
	storage.CreateColumn(table.BuildColumnOf[components.Location](columns.NewBit(components.Location{})))
	storage.CreateColumn(table.BuildColumnOf[components.Locked](columns.NewBit(components.Locked{})))
	storage.CreateColumn(table.BuildColumnOf[components.MasterQuest](columns.NewBit(components.MasterQuest{})))
	storage.CreateColumn(table.BuildColumnOf[components.Placeable](columns.NewBit(components.Placeable{})))
	storage.CreateColumn(table.BuildColumnOf[components.Pots](columns.NewBit(components.Pots{})))
	storage.CreateColumn(table.BuildColumnOf[components.Refill](columns.NewBit(components.Refill{})))
	storage.CreateColumn(table.BuildColumnOf[components.RupeeTowers](columns.NewBit(components.RupeeTowers{})))
	storage.CreateColumn(table.BuildColumnOf[components.SmallCrates](columns.NewBit(components.SmallCrates{})))
	storage.CreateColumn(table.BuildColumnOf[components.Token](columns.NewBit(components.Token{})))

	return nil
}
