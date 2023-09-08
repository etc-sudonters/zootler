package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/etc-sudonters/zootler/internal/datastructures/graph"
	"github.com/etc-sudonters/zootler/internal/rules"
	"github.com/etc-sudonters/zootler/pkg/entity"
	"github.com/etc-sudonters/zootler/pkg/entity/hashpool"
	"github.com/etc-sudonters/zootler/pkg/world"
)

func main() {
	args := os.Args

	var dir string

	if len(args) > 1 {
		dir = args[1]
	} else {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var allErrs []error
	b := &WorldBuilder{
		G: graph.Builder{
			graph.WithCapacity(512),
		},
		P: hashpool.New(),
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fp := path.Join(dir, entry.Name())
		locs, err := rules.ReadLogicFile(fp)
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("failed to read %s: %w", fp, err))
			continue
		}

		for _, loc := range locs {
			fmt.Printf("Accepting %s", loc)
			b.Accept(loc)
		}
	}

	if allErrs != nil {
		fmt.Println(errors.Join(allErrs...))
		os.Exit(1)
	}
}

func accept(w *world.World, loc rules.RawLogicLocation) error {

}

type WorldBuilder struct {
	G graph.Builder
	P entity.Pool
}

func (w *WorldBuilder) Accept(loc rules.RawLogicLocation) error {
	return nil
}
