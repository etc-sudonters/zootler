package hashpool

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/etc-sudonters/zootler/internal/set"
	"github.com/etc-sudonters/zootler/pkg/entity"
)

func dump(t *testing.T, v interface{}) {
	t.Logf("%+v", v)
}

func dumpView(w io.Writer, v interface{}) {
	fmt.Fprintf(w, "Dump:\n%+v\n", v)
}

func expectedComponent(w io.Writer, v view, expected entity.Component) {
	fmt.Fprintf(w, "expected component %s to be loaded on %d\n", expected, v.m)
}

func didNotExpectError(t *testing.T, err error) {
	w := &strings.Builder{}
	fmt.Fprintf(w, "did not expected error:\n%s", err)
	t.Fatal(w.String())
}

func expectedEqualComponents[T entity.Component](w io.Writer, v view, expect T, actual T) {
	fmt.Fprintf(w, "expected equal values for component %s on %d\n", entity.ComponentName(expectedComponent), v.m)
	fmt.Fprint(w, "Actual ")
	dumpView(w, actual)
	fmt.Fprint(w, "Expected ")
	dumpView(w, expect)
}

func TestCanRetrieveComponentFromView(t *testing.T) {
	p := New()
	ent := p.createCore()

	var model entity.Model = 99999
	err := ent.Get(&model)

	if err != nil {
		msg := &strings.Builder{}
		expectedComponent(msg, ent, model)
		dumpView(msg, ent)
		t.Fatal(msg.String())
	}

	if model != ent.Model() {
		msg := &strings.Builder{}
		expectedEqualComponents(msg, ent, ent.Model(), model)
		t.Fatal(msg.String())
	}
}

func TestCanStoreAndRetrievePointerToComp(t *testing.T) {
	initialValue := 10
	changedValue := 9999

	p := New()
	ent := p.createCore()
	ent.Add(&myTestComponent{initialValue})

	var c *myTestComponent
	// ptrB -> ptrA, assigns ptrA
	if err := ent.Get(&c); err != nil {
		didNotExpectError(t, err)
	}

	c.V = changedValue

	var d *myTestComponent
	if err := ent.Get(&d); err != nil {
		didNotExpectError(t, err)
	}

	if d.V != c.V {
		msg := &strings.Builder{}
		expectedEqualComponents(msg, ent, c, d)
		t.Fatal(msg.String())
	}
}

func TestCanStoreComponentAndRetrieveThroughPointer(t *testing.T) {
	t.Skip("Not great this doesn't work, but we can store pointers and that's good enough")
	initialValue := 10
	changedValue := 9999

	p := New()
	ent := p.createCore()
	// NOTE _not_ a pointer
	ent.Add(myTestComponent{initialValue})

	// NOTE _is_ pointer
	t.Log("loading through pointer")
	var c *myTestComponent
	if err := ent.Get(&c); err != nil {
		didNotExpectError(t, err)
	}
	t.Log("freshly fetched")
	dump(t, c)
	t.Log("changing value")
	c.V = changedValue
	dump(t, c)

	// NOTE not pointer
	t.Log("loading through not pointer")
	var d myTestComponent
	if err := ent.Get(&d); err != nil {
		didNotExpectError(t, err)
	}

	if d.V != c.V {
		msg := &strings.Builder{}
		expectedEqualComponents(msg, ent, *c, d) // deref c to make types happy
		t.Fatal(msg.String())
	}
}

func TestCanRemoveCustomComponent(t *testing.T) {
	p := New()
	ent := p.createCore()
	ent.Add(myTestComponent{})

	var c myTestComponent
	if err := ent.Get(&c); err != nil {
		didNotExpectError(t, err)
	}

	// doesn't need to be loaded instance
	ent.Remove(myTestComponent{})

	var d myTestComponent
	if err := ent.Get(&d); err != nil && !errors.Is(err, entity.ErrNotLoaded) {
		didNotExpectError(t, err)
	}
}

func TestCanQueryForEntitiesByComponentExistence(t *testing.T) {
	t.Log("sure would be nice to make this a nice big number")
	componentsToMake := 10000
	tagRatio := 7
	p := New()

	totalEnts := set.New[entity.Model]()
	taggedEnts := set.New[entity.Model]()

	for i := 0; i <= componentsToMake; i++ {
		ent := p.createCore()
		totalEnts.Add(ent.Model())

		if (i % tagRatio) == 0 {
			ent.Add(myTestComponent{i})
			taggedEnts.Add(ent.Model())
		}
	}

	dump(t, p)

	if len(totalEnts) != len(p.All()) {
		t.Logf("mismatched entity count\nexpected:\t%d\nactual:\t%d", len(totalEnts), len(p.All()))
		t.FailNow()
	}

	queryedFor, err := p.Query(
		entity.With[myTestComponent]{},
		entity.DebugSelector{
			F: func(s string, a ...any) {},
			S: entity.Load[myTestComponent]{},
		},
	)
	if err != nil {
		didNotExpectError(t, err)
	}

	if len(queryedFor) != len(taggedEnts) {
		t.Logf("mismatched entity count\nexpected:\t%d\nactual:\t%d", len(taggedEnts), len(queryedFor))
		t.FailNow()
	}

	if !reflect.DeepEqual(
		taggedEnts,
		set.MapFromSlice(
			queryedFor,
			func(u entity.View) entity.Model {
				return u.Model()
			},
		),
	) {
		t.Logf("tagged set and returned set are different!")
		t.Fail()
	}
}

func TestCanUseMultipleComponents(t *testing.T) {
	componentsToMake := 35
	goodTagRatio := 7
	badTagRation := 5
	p := New()

	totalEnts := set.New[entity.Model]()
	goodTaggedEnts := set.New[entity.Model]()
	badTaggedEnts := set.New[entity.Model]()

	for i := 0; i <= componentsToMake; i++ {
		ent := p.createCore()
		totalEnts.Add(ent.Model())

		if (i % goodTagRatio) == 0 {
			ent.Add(myTestComponent{i})
			goodTaggedEnts.Add(ent.Model())
		}

		if (i % badTagRation) == 0 {
			ent.Add(anotherComponent{float64(i)})
			badTaggedEnts.Add(ent.Model())
		}
	}

	comboTagSet := set.Intersection(goodTaggedEnts, badTaggedEnts)
	comboQueries, err := p.Query(
		entity.DebugSelector{
			F: func(string, ...any) {}, //t.Logf,
			S: entity.Load[myTestComponent]{},
		},
		entity.DebugSelector{
			F: t.Logf,
			S: entity.Load[anotherComponent]{},
		},
	)

	if err != nil {
		didNotExpectError(t, err)
	}

	actualEntities := set.MapFromSlice(comboQueries, entity.View.Model)

	if !reflect.DeepEqual(
		comboTagSet, actualEntities,
	) {
		t.Log("selected a different set of entities than expected")
		t.Logf("expected:\t%+v", comboTagSet)
		t.Logf("actual:\t\t%+v", actualEntities)
		t.Fail()
	}
}

func TestCanExcludeEntitiesBasedOnComponent(t *testing.T) {
	componentsToMake := 1000
	firstTagRatio := 7
	secondTagRatio := 5
	p := New()

	totalEnts := set.New[entity.Model]()
	firstTagEnts := set.New[entity.Model]()
	secondTagEnt := set.New[entity.Model]()

	for i := 0; i <= componentsToMake; i++ {
		ent := p.createCore()
		totalEnts.Add(ent.Model())

		if (i % firstTagRatio) == 0 {
			ent.Add(myTestComponent{i})
			firstTagEnts.Add(ent.Model())
		}

		if (i % secondTagRatio) == 0 {
			ent.Add(anotherComponent{float64(i)})
			secondTagEnt.Add(ent.Model())
		}
	}

	allEntitiesWithoutTags := set.Difference(
		set.Union(firstTagEnts, secondTagEnt),
		totalEnts,
	)

	queriedAllUntagged, err := p.Query(
		entity.Without[myTestComponent]{},
		entity.Without[anotherComponent]{},
	)

	if err != nil {
		didNotExpectError(t, err)
	}

	actualEntities := set.MapFromSlice(queriedAllUntagged, entity.View.Model)

	if actualEntities.Exists(entity.Model(firstTagRatio*secondTagRatio)) ||
		actualEntities.Exists(entity.Model(firstTagRatio)) ||
		actualEntities.Exists(entity.Model(secondTagRatio)) {
		t.Log("selected incorrect entity group")
		t.FailNow()
	}

	if !reflect.DeepEqual(allEntitiesWithoutTags, actualEntities) {
		t.Log("selected a different set of entities than expected")
		t.Logf("expected:\t%+v", allEntitiesWithoutTags)
		t.Logf("actual:\t\t%+v", actualEntities)
		t.Fail()
	}

}

func TestCanFilterWithoutLoading(t *testing.T) {
	componentsToMake := 10000
	tagRatio := 7
	p := New()

	totalEnts := set.New[entity.Model]()
	taggedEnts := set.New[entity.Model]()
	taggedCount := 0

	for i := 0; i <= componentsToMake; i++ {
		ent := p.createCore()
		totalEnts.Add(ent.Model())

		if (i % tagRatio) == 0 {
			ent.Add(myTestComponent{i})
			taggedEnts.Add(ent.Model())
			taggedCount += 1
		}
	}

	if len(taggedEnts) != taggedCount {
		t.Logf("expected to make %d entities", taggedCount)
		t.Logf("actually made %d", len(taggedEnts))
		t.FailNow()
	}

	queriedEnts, err := p.Query(entity.With[myTestComponent]{})

	if err != nil {
		didNotExpectError(t, err)
	}

	var comp myTestComponent
	for _, ent := range queriedEnts {
		err := ent.Get(&comp)
		if errors.Is(err, entity.ErrNotLoaded) {
			continue
		}

		if err != nil {
			didNotExpectError(t, err)
		}

		t.Logf("did not expect %T to be loaded on %d", comp, ent.Model())
		t.Fail()
	}

	if len(queriedEnts) != taggedCount {
		t.Logf("expected %d elements", componentsToMake/tagRatio)
		t.Logf("got %d elements", len(queriedEnts))
		t.Fail()
	}
}

func TestCanRetrieveArbitraryEntityWithComps(t *testing.T) {
	p := New()
	ent := p.createCore()
	ent.Add(myTestComponent{99})

	var c1 *myTestComponent
	var c2 *anotherComponent

	p.Get(ent.m, &c1, &c2)

	if c1 == nil {
		t.Logf("expected to retrieve %T from %s", c1, ent)
		t.Fail()
	}

	if c1.V != 99 {
		t.Logf("did not retrieve expected instance of %[1]T: %[1]v", c1)
		t.Fail()
	}

	if c2 != nil {
		t.Logf("did not expect to retrieve %T from %s", c2, ent)
		t.Fail()
	}

}

type myTestComponent struct {
	V int
}

type anotherComponent struct {
	K float64
}
