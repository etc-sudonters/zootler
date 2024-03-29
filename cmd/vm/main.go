package main

import (
	"fmt"
	"sudonters/zootler/pkg/rules/vm"
)

func main() {
	var c vm.Chunk
	c.Name = "test"
	c.Code = []vm.Op{
		vm.OP_CONSTANT, c.AddConstant(vm.Value{
			Kind:  vm.ValNum,
			Value: 3.14,
		}),
		vm.OP_SELECT, c.AddConstant(vm.Value{
			Kind:  vm.ValNum,
			Value: 18,
		}),
		vm.OP_WITH, c.AddConstant(vm.Value{
			Kind:  vm.ValNum,
			Value: 8237,
		}),
		vm.OP_WITHOUT, c.AddConstant(vm.Value{
			Kind:  vm.ValNum,
			Value: 1 << 54,
		}),
		vm.OP_RETURN,
	}
	fmt.Print(c.String())
}
