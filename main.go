package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/etc-sudonters/zootler/internal/rules"
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

		if err = rules.LexAllLocationRules(locs); err != nil {
			allErrs = append(
				allErrs,
				fmt.Errorf("failures in %s: %w", fp, err),
			)
		}
	}

	if allErrs != nil {
		fmt.Println(errors.Join(allErrs...))
		os.Exit(1)
	}
}
