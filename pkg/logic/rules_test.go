package logic

import (
	"testing"

	"github.com/etc-sudonters/zootler/pkg/entity/hashpool"
)

func TestHasQuantityOf(t *testing.T) {
	pool := hashpool.New()
	desired := Name("compilers")

	others := []Name{"lexer", "parser", "interpreter", "linter"}

	for i := 0; i < 100; i++ {
		ent, _ := pool.Create()
		ent.Add(Token{})
		ent.Add(others[i%len(others)])
		ent.Add(Collected{})
	}

	for i := 0; i < 4; i++ {
		ent, _ := pool.Create()
		ent.Add(Token{})
		ent.Add(desired)
		ent.Add(Collected{})
	}

	rule := hasQuantityOf{desired, 4}

	fulfilled, _ := rule.Fulfill(pool)

	if !fulfilled {
		t.Logf("expected 4 %s but did not find them", desired)
		t.Fail()
	}
}
