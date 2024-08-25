package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/app"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/preprocessor"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/slipup"
)

var builtInNames []string = []string{
	"has", "at_night", "at_day", "at_dampe_time",
	"has_bottle", "had_night_start", "has_all_notes_for_song",
}

type WorldLoader struct {
	Path, Helpers string
	Settings      map[string]any
}

func (w WorldLoader) Setup(z *app.Zootlr) error {
	ctx := z.Ctx()
	eng := z.Engine()
	settings := app.GetResource[settings.ZootrSettings](z)

	assistant, createAssistantErr := w.assistant(ctx, eng, settings.Res)
	if createAssistantErr != nil {
		return slipup.Describe(createAssistantErr, "while create world loader assistant")
	}

	if dirLoadErr := w.loadDirectory(ctx, assistant); dirLoadErr != nil {
		return slipup.Describe(dirLoadErr, "while loading logic directory")
	}

	if processErr := w.processLocations(ctx, assistant); processErr != nil {
		return slipup.Describe(processErr, "while preprocessing locations")
	}

	return nil
}

func (w WorldLoader) processLocations(_ context.Context, a *assistant) error {
	pre := a.preprocessor()
	var errs []error

	for loc := range a.locs.loaded {
		processed, err := pre.Process(string(loc.name), loc.rules)
		if err != nil {
			err = slipup.Describef(err, "while preprocessing location '%s'", loc.name)
			errs = append(errs, err)
			continue
		}

		loc.rules = processed
	}

	return errors.Join(errs...)
}

func (w WorldLoader) loadDirectory(ctx context.Context, a *assistant) error {
	dontio.WriteLineOut(ctx, "reading dir  '%s'", w.Path)
	directory, dirErr := os.ReadDir(w.Path)
	if dirErr != nil {
		return slipup.Describe(dirErr, w.Path)
	}
	for _, entry := range directory {
		if !internal.IsFile(entry) {
			continue
		}

		path := path.Join(w.Path, entry.Name())
		if logicFileErr := a.loadLogicFile(path); logicFileErr != nil {
			return slipup.Describe(logicFileErr, path)
		}
	}

	return nil
}

func (w WorldLoader) assistant(ctx context.Context, eng query.Engine, settings settings.ZootrSettings) (*assistant, error) {
	a := new(assistant)
	if locs, err := createLocationTable(ctx, eng); err != nil {
		return nil, slipup.Describe(err, "while building location table")
	} else {
		a.locs = locs
	}

	if ft, ftErr := createFunctionTable(ctx, w.Helpers); ftErr != nil {
		return nil, slipup.Describe(ftErr, "while create function table")
	} else {
		a.ft = ft
	}

	if vals, valErr := createValuesTable(ctx, eng, settings); valErr != nil {
		return nil, slipup.Describe(valErr, "while creating values table")
	} else {
		a.vals = vals
	}
	return a, nil
}

type assistant struct {
	ft   preprocessor.FunctionTable
	locs *LocationTable
	vals preprocessor.ValuesTable
}

func (a assistant) preprocessor() *preprocessor.P {
	return &preprocessor.P{
		Delayed:   make(preprocessor.DelayedRules, 256),
		Functions: a.ft,
		Env:       a.vals,
	}
}

func (a assistant) loadLogicFile(path string) error {
	if strings.Contains(path, "mq") {
		return nil
	}

	rawLocs, readErr := internal.ReadJsonFileAs[[]logicLocation](path)
	if readErr != nil {
		return slipup.Describe(readErr, path)
	}

	for _, raw := range rawLocs {
		here, buildErr := a.locs.Build(components.Name(raw.Name))
		if buildErr != nil {
			return slipup.Describef(buildErr, "building %s", raw.Name)
		}

		values := table.Values{
			components.HintRegion{
				Name: raw.HintRegion,
				Alt:  raw.HintRegionAlt,
			},
		}

		if raw.BossRoom {
			values = append(values, components.BossRoom{})
		}

		if raw.Dungeon != "" {
			values = append(values, components.Dungeon(raw.Dungeon))
		}

		if raw.DoesTimePass {
			values = append(values, components.TimePasses{})
		}

		if raw.Savewarp != "" {
			values = append(values, components.SavewarpName(raw.Savewarp))
		}

		if err := here.add(values); err != nil {
			return slipup.Describef(err, "while populating %s", raw.Name)
		}

		for destination, rule := range raw.Exits {
			linkErr := here.linkTo(destination, rule, components.ExitEdge{})
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, destination)
			}
		}

		for check, rule := range raw.Locations {
			linkErr := here.linkTo(check, rule, components.CheckEdge{})
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, check)
			}
		}

		for event, rule := range raw.Events {
			linkErr := here.linkTo(event, rule, components.EventEdge{})
			// TODO: create a token for the event as well
			if linkErr != nil {
				return slipup.Describef(linkErr, "while linking '%s -> %s'", here.name, event)
			}
		}

	}

	return nil

}

type logicLocation struct {
	BossRoom      bool              `json:"is_boss_room"`
	DoesTimePass  bool              `json:"time_passes"`
	Dungeon       string            `json:"dungeon"`
	HintRegion    string            `json:"hint"`     // hint zone
	HintRegionAlt string            `json:"alt_hint"` // more specific location -- restricted to gannon's chambers?
	Name          string            `json:"region_name"`
	Savewarp      string            `json:"savewarp"` // effectively another exit
	Scene         string            `json:"scene"`    // similar to hint?
	Events        map[string]string `json:"events"`
	Exits         map[string]string `json:"exits"`
	Locations     map[string]string `json:"locations"`
}

type locationBuilder struct {
	parent *LocationTable
	id     table.RowId
	name   components.Name
	rules  map[string]parser.Expression
}

type parsedRule struct {
	into string
	ast  parser.Expression
}

func (l locationBuilder) add(v table.Values) error {
	return l.parent.queries.SetValues(l.id, v)
}

func (l *locationBuilder) linkTo(name, rule string, vs ...table.Value) error {
	edgeName := fmt.Sprintf("%s -> %s", l.name, name)

	ast, parseErr := parser.Parse(rule)
	if parseErr != nil {
		return slipup.Describef(parseErr, "while parsing rule '%s'", rule)
	}

	destination, linkErr := l.parent.Build(components.Name(name))
	if linkErr != nil {
		return slipup.Describef(linkErr, "while linking '%s'", edgeName)
	}

	edge, edgeCreateErr := l.parent.queries.InsertRow(components.Name(edgeName), components.Edge{
		Origin: entity.Model(l.id),
		Dest:   entity.Model(destination.id),
	})

	if edgeCreateErr != nil {
		return slipup.Describef(edgeCreateErr, "while creating edge %s", edgeName)
	}

	if len(vs) != 0 {
		if additionalErr := l.parent.queries.SetValues(edge, table.Values(vs)); additionalErr != nil {
			return slipup.Describef(additionalErr, "while customizing '%s'", edgeName)
		}
	}

	l.rules[name] = ast
	return nil
}

func createLocationTable(_ context.Context, e query.Engine) (*LocationTable, error) {
	q := e.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](e))
	q.Exists(query.MustAsColumnId[components.Location](e))
	rows, qErr := e.Retrieve(q)
	if qErr != nil {
		return nil, qErr
	}

	tbl := new(LocationTable)
	tbl.queries = e
	tbl.locations = make(map[internal.NormalizedStr]*locationBuilder, rows.Len())

	for id, tup := range rows.All {
		name, ok := tup.Values[0].(components.Name)
		if !ok {
			return nil, fmt.Errorf("could not cast row %d %+v to 'components.Name'", id, tup)
		}

		tbl.locations[internal.Normalize(name)] = tbl.builder(name, id)
	}

	return tbl, nil
}

type LocationTable struct {
	queries   query.Engine
	locations map[internal.NormalizedStr]*locationBuilder
}

func (l *LocationTable) Build(name components.Name) (*locationBuilder, error) {
	normalized := internal.Normalize(name)
	if existing, ok := l.locations[normalized]; ok {
		return existing, nil
	}

	row, insertErr := l.queries.InsertRow(name, components.Location{})
	if insertErr != nil {
		return nil, insertErr
	}
	b := l.builder(name, row)
	l.locations[normalized] = b
	return b, nil
}

func (l *LocationTable) builder(name components.Name, row table.RowId) *locationBuilder {
	b := new(locationBuilder)
	b.id = row
	b.name = name
	b.parent = l
	b.rules = make(map[string]parser.Expression, 16)
	return b
}

func (l *LocationTable) loaded(yield func(*locationBuilder) bool) {
	for _, b := range l.locations {
		if !yield(b) {
			break
		}
	}
}

func createFunctionTable(_ context.Context, helperPath string) (preprocessor.FunctionTable, error) {
	helpers, helperReadErr := internal.ReadJsonFileStringMap(helperPath)
	if helperReadErr != nil {
		return nil, slipup.Describef(helperReadErr, "while reading '%s'", helperPath)
	}

	table := make(preprocessor.FunctionTable, len(helpers))
	var errs []error

	for decl, body := range helpers {
		f, parseErr := parser.ParseFunctionDecl(decl, body)
		if parseErr != nil {
			errs = append(errs, slipup.Describef(parseErr, "while parsing function decl '%s'", decl))
			continue
		}

		table[f.Identifier] = f
	}

	for _, name := range builtInNames {
		table.BuiltIn(name)
	}

	return table, errors.Join(errs...)
}

func createValuesTable(_ context.Context, e query.Engine, s settings.ZootrSettings) (preprocessor.ValuesTable, error) {
	tbl := make(preprocessor.ValuesTable, 512)

	if err := func() error {
		q := e.CreateQuery()
		q.Load(query.MustAsColumnId[components.Name](e))
		q.Exists(query.MustAsColumnId[components.CollectableGameToken](e))
		q.Exists(query.MustAsColumnId[components.Advancement](e))
		rows, qErr := e.Retrieve(q)
		if qErr != nil {
			return qErr
		}

		for id, tup := range rows.All {
			name, ok := tup.Values[0].(components.Name)
			if !ok {
				return fmt.Errorf("could not cast row %d %+v to 'components.Name'", id, tup)
			}
			normaled := internal.Normalize(name)
			tbl[normaled] = parser.TokenLiteral(uint64(id))
		}
		return nil
	}(); err != nil {
		return nil, err
	}

	tbl["shuffleindividualocarinanotes"] = parser.BoolLiteral(s.Shuffling.OcarinaNotes)
	tbl["fixbrokendrops"] = parser.BoolLiteral(s.FixBrokenDrops)
	tbl["plantbeans"] = parser.BoolLiteral(s.Starting.Beans)
	tbl["bigpoecount"] = parser.NumberLiteral(s.Minigames.BigPoeCount)
	tbl["shuffleexpensivemerchants"] = parser.BoolLiteral(s.Shuffling.ExpensiveMerchants)
	tbl["freebombchudrops"] = parser.BoolLiteral(s.FreeBombchuDrops)
	tbl["shuffleemptypots"] = parser.BoolLiteral(settings.HasFlag(s.Shuffling.Pots, settings.ShuffleEmptyPots))
	tbl["zorafountain"] = parser.StringLiteral(s.Locations.ZoraFountain.String())
	tbl["damagemultiplier"] = parser.StringLiteral(s.Damage.Multiplier.String())
	tbl["deadlybonks"] = parser.StringLiteral(s.Damage.Bonk.String())
	tbl["keysanity"] = parser.BoolLiteral(false)
	tbl["shufflepots"] = parser.StringLiteral("off")
	tbl["gerudofortress"] = parser.StringLiteral("closed")
	tbl["freescarecrow"] = parser.BoolLiteral(s.Starting.Scarecrow)
	tbl["shuffleoverworldentrances"] = parser.BoolLiteral(s.Entrances.Overworld)
	tbl["shuffleinteriorentrances"] = parser.BoolLiteral(s.Entrances.Interior != settings.InteriorShuffleOff)
	tbl["shufflesilverrupees"] = parser.BoolLiteral(s.KeyShuffle.SilverRupees&settings.KeysVanilla != settings.KeysVanilla)
	tbl["bridge"] = parser.StringLiteral("open")
	tbl["shuffletcgkeys"] = parser.BoolLiteral(s.KeyShuffle.ChestGameKeys&settings.KeysVanilla != settings.KeysVanilla)
	tbl["completemaskquest"] = parser.BoolLiteral(settings.HasFlag(s.Trades.Child, settings.ChildTradeComplete))
	tbl["openkakariko"] = parser.StringLiteral("open")
	tbl["adulttradeshuffle"] = parser.BoolLiteral(true)
	tbl["selectedadulttradeitem"] = parser.StringLiteral("Claim Check")
	tbl["skiprewardfromrauru"] = parser.BoolLiteral(s.Starting.RauruReward)
	tbl["skipchildzelda"] = parser.BoolLiteral(s.Locations.SkipChildZelda)
	tbl["entranceshuffle"] = parser.BoolLiteral(false)
	tbl["shuffledungeonentrances"] = parser.BoolLiteral(false)
	tbl["openforest"] = parser.StringLiteral("closed_deku")
	tbl["disabletraderevert"] = parser.BoolLiteral(true)
	tbl["opendooroftime"] = parser.BoolLiteral(s.Locations.OpenDoorOfTime)
	tbl["foresttempleamyandmeg"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.ForestTemplePoes, settings.ForestTempleAmyMeg))
	tbl["foresttemplejoandbeth"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.ForestTemplePoes, settings.ForestTempleJoBeth))
	tbl["chickencount"] = parser.NumberLiteral(s.Minigames.KakChickens)
	tbl["shufflescrubs"] = parser.BoolLiteral(s.Shuffling.Scrubs != settings.ShuffleScrubsOff && s.Shuffling.Scrubs != settings.ShuffleScrubsUpgradeOnly)

	trials := [6]internal.NormalizedStr{"forest", "fire", "water", "shadow", "spirit", "light"}
	switch s.Dungeons.Trials {
	case settings.TrialsEnabledNone:
		for _, t := range trials {
			tbl["skippedtrials"+t] = parser.BoolLiteral(true)
		}
		break
	case settings.TrialsEnabledAll:
		for _, t := range trials {
			tbl["skippedtrials"+t] = parser.BoolLiteral(false)
		}
		break
	default:
		tbl["skipppedtrialsforest"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledForest))
		tbl["skipppedtrialsfire"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledFire))
		tbl["skipppedtrialswater"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledWater))
		tbl["skipppedtrialsshadow"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledShadow))
		tbl["skipppedtrialsspirit"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledSpirit))
		tbl["skipppedtrialslight"] = parser.BoolLiteral(settings.HasFlag(s.Dungeons.Trials, settings.TrialsEnabledLight))
	}

	return tbl, nil
}
