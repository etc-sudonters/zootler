package main

import (
	"fmt"

	"github.com/etc-sudonters/zootler/entity"
	"github.com/etc-sudonters/zootler/entity/hashpool"
)

func main() {
	_, _ = hashpool.New[entity.Model]()
	fmt.Println("vim-go")
}
