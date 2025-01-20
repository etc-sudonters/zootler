package main

import (
	"fmt"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/internal/table/columns"
)

func displaycolstats(tbl *table.Table) {
	fmt.Printf("%6d rows\n", len(tbl.Rows))
	fmt.Printf("%4d columns\n", len(tbl.Cols))
	fmt.Println()
	for _, c := range tbl.Cols {
		col := c.Column()
		fmt.Printf("%6d ", col.Len())
		slc, isslc := col.(*columns.Slice)
		if isslc {
			fmt.Printf("%6d ", slc.Capacity())
		} else {
			fmt.Print("       ")
		}
		fmt.Printf("%s\n", c.Type().Name())
	}
}
