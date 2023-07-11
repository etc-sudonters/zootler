package hashpool

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/etc-sudonters/zootler/entity"
	"github.com/etc-sudonters/zootler/set"
)

func newTestingPool(t *testing.T) *Pool {
	p, err := New()
	if err != nil {
		t.Fatalf("could not initialize pool: %s", err)
	}
	p.debug = t.Logf
	return p
}

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
	p := newTestingPool(t)
	ent := p.createEasy()

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

	p := newTestingPool(t)
	ent := p.createEasy()
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

	p := newTestingPool(t)
	ent := p.createEasy()
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
	p := newTestingPool(t)
	ent := p.createEasy()
	ent.Add(myTestComponent{})

	var c myTestComponent
	if err := ent.Get(&c); err != nil {
		didNotExpectError(t, err)
	}

	// doesn't need to be loaded instance
	ent.Remove(myTestComponent{})

	var d myTestComponent
	if err := ent.Get(&d); err != nil && !errors.Is(err, ErrNotLoaded) {
		didNotExpectError(t, err)
	}
}

func TestCanQueryForEntitiesByComponentExistence(t *testing.T) {
	t.Log("sure would be nice to make this a nice big number")
	componentsToMake := 10000
	tagRatio := 7
	p := newTestingPool(t)

	totalEnts := set.New[entity.Model]()
	taggedEnts := set.New[entity.Model]()

	for i := 0; i <= componentsToMake; i++ {
		ent := p.createEasy()
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

	t.Logf("issusing query for Include[myTestComponent]")
	queryedFor, err := p.Query(
		myTestComponent{},
		entity.DebugSelector{
			Debug:    t.Logf,
			Selector: entity.Include[myTestComponent]{},
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

type myTestComponent struct {
	V int
}
