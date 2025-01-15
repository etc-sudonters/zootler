package z2

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/symbols"

	"github.com/etc-sudonters/substrate/slipup"
)

type attachments struct {
	v table.Values
}

func (this *attachments) add(v table.Value) {
	this.v = append(this.v, v)
}

func AliasTokens(symbols *symbols.Table, funcs *ast.PartialFunctionTable, names []string) {
	for _, name := range names {
		escaped := escape(name)
		if _, exists := funcs.Get(escaped); exists {
			continue
		}
		if _, exists := funcs.Get(name); exists {
			continue
		}
		symbol := symbols.LookUpByName(name)
		symbols.Alias(symbol, escaped)
	}
}

var escaping = regexp.MustCompile("['()[\\]-]")

func escape(name string) string {
	name = escaping.ReplaceAllLiteralString(name, "")
	return strings.ReplaceAll(name, " ", "_")
}

func ReadHelpers(path string) map[string]string {
	contents, err := internal.ReadJsonFileStringMap(path)
	if err != nil {
		panic(err)
	}
	return contents
}

func LoadRegionsFrom(dir string) ([]region, error) {
	var regions []region
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return slipup.Describe(err, "logic directory walk called with err")
		}

		info, err := entry.Info()
		if err != nil || info.Mode() != (^fs.ModeType)&info.Mode() {
			// either we couldn't get the info, which doesn't bode well
			// or it's some kind of not file thing which we also don't want
			return nil
		}

		if ext := filepath.Ext(path); ext != ".json" {
			return nil
		}

		these, readErr := internal.ReadJsonFileAs[[]region](path)
		if readErr != nil {
			return slipup.Describef(readErr, "while reading file '%s'", path)
		}

		regions = slices.Concat(regions, these)
		return nil
	})

	return regions, err
}

func LoadTokensFrom(path string) ([]token, error) {
	return internal.ReadJsonFileAs[[]token](path)
}

type region struct {
	Events     map[string]string `json:"events"`
	Exits      map[string]string `json:"exits"`
	Locations  map[string]string `json:"locations"`
	RegionName string            `json:"region_name"`
	AltHint    string            `json:"alt_hint"`
	Hint       string            `json:"hint"`
	Dungeon    string            `json:"dungeon"`
	IsBossRoom bool              `json:"is_boss_room"`
	Savewarp   string            `json:"savewarp"`
	Scene      string            `json:"scene"`
	TimePasses bool              `json:"time_passes"`
}

type token struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Advancement bool                   `json:"advancement"`
	Priority    bool                   `json:"priority"`
	Special     map[string]interface{} `json:"special"`
}
