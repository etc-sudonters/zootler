package magicbeanvm

import (
	"fmt"
	"sudonters/zootler/midologic/objects"
)

type TranslationLayer struct {
	byaddr map[objects.Ptr]CollectionId
	byname map[string]CollectionId
}

func (this TranslationLayer) FromName(name string) CollectionId {
	id, exists := this.byname[name]
	if !exists {
		panic(fmt.Errorf("%q not declared in address table", name))
	}
	return id
}

func (this TranslationLayer) FromPtr(ptr objects.Ptr) CollectionId {
	id, exists := this.byaddr[ptr]
	if !exists {
		panic(fmt.Errorf("bad token pointer deref: %s", ptr))
	}

	return id
}

type TranslationLayerBuilder struct {
	*TranslationLayer
}

func (this TranslationLayerBuilder) Declare(name string, id CollectionId) {
	if already, exists := this.byname[name]; exists {
		if already == id {
			return
		}
		panic(fmt.Errorf("%q reassigned collection id from %d to %d", name, already, id))
	}

	this.byname[name] = id
}

func (this TranslationLayerBuilder) PointTo(name string, tag objects.PtrTag) objects.Ptr {
	id, exists := this.byname[name]
	if !exists {
		panic(fmt.Errorf("%q is not declared in address table", name))
	}

	ptr := objects.Pointer(objects.OpaquePointer(id), tag)
	this.byaddr[ptr] = id
	return ptr
}
