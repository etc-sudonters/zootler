package main

import (
	"fmt"

	"github.com/etc-sudonters/rando/entity"
	"github.com/etc-sudonters/rando/entity/hashpool"
)

func main() {
	var p entity.Pool = hashpool.New()
	fmt.Println("vim-go")
}
