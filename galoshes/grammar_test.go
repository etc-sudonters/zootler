package galoshes

import (
	"testing"

	"github.com/etc-sudonters/substrate/peruse"
)

func TestCanParseFind(t *testing.T) {
	Repr = PrettyPrint
	findScript := `find [ $place-id $place-name ]
where [
    [ :holding $place-id $token-id ]
    [ :named $place-id $place-name ]
    [ :named $token-id "Kokri Sword" ]
    [ $token-id tokens 1 ] 
]
rules [
    [:holding $place $token] [ [$place world/placement/holds $token] ]
    [:named $id $name] [ [$id names $name] ]
]
`
	expectedFinding := 2
	expectedConstraints := 4
	expectedDerivations := 2

	lexer := NewLexer(findScript)
	parser := peruse.NewParser(grammar, lexer)

	result, parseErr := parser.ParseAt(peruse.LOWEST)
	t.Logf("lexer state: %#v", lexer)
	t.Logf("parser state: %#v", parser)
	if parseErr != nil {
		t.Fatal(parseErr.Error())
	}

	find, isFind := assertIsType[Find](t, result)

	if !isFind {
		t.FailNow()
	}

	t.Logf("%s", find)

	if len(find.Finding) != expectedFinding {
		t.Fatalf("expected to parse %d variables, parsed %d", expectedFinding, len(find.Finding))
	}

	if find.Finding[0] != Variable("place-id") {
		t.Logf("expected index 0 to be place-id but was %s", find.Finding[0])
		t.Fail()
	}

	if find.Finding[1] != Variable("place-name") {
		t.Logf("expected index 1 to be place-id but was %s", find.Finding[1])
		t.Fail()
	}

	if len(find.Constraints) != expectedConstraints {
		t.Logf("expected to parse %d constraints, parsed %d", expectedConstraints, len(find.Constraints))
		t.Fail()
	}

	for i := range len(find.Constraints) {
		switch i {
		case 0:
			assertIsInvoke(t, find.Constraints[0], Invoke("holding", VarOf[Literal]("place-id"), VarOf[Literal]("token-id")))
		case 1:
			assertIsInvoke(t, find.Constraints[1], Invoke("named", VarOf[Literal]("place-id"), VarOf[Literal]("place-name")))
		case 2:
			assertIsInvoke(t, find.Constraints[2], Invoke("named", VarOf[Literal]("token-id"), Part(LiteralString("Kokri Sword"))))
		case 3:
			assertIsTriplet(t, find.Constraints[3], Triplet{
				Id:    VarOf[Number]("token-id"),
				Attr:  Part(Attribute("tokens")),
				Value: Part(LiteralNumber(1)),
			})
		}
	}

	if len(find.Derivations) != expectedDerivations {
		t.Logf("expected to parse %d derivations, parsed %d", expectedDerivations, len(find.Derivations))
		t.Fail()
	}

	for i := range find.Derivations {
		derivation, isD := assertIsType[DerivationDecl](t, find.Derivations[i])
		if !isD {
			continue
		}
		switch i {
		case 0:
			if assertIsDecl(t, derivation, "holding", "place", "token") {
				if len(derivation.Constraints) != 1 {
					t.Logf("expected 1 constraint, found %d", len(derivation.Constraints))
					t.Fail()
				} else {
					assertIsTriplet(t, derivation.Constraints[0], Triplet{
						Id:    VarOf[Number]("place"),
						Attr:  Part[Attribute]("world/placement/holds"),
						Value: VarOf[Literal]("token"),
					})
				}
			}
		case 1:
			if assertIsDecl(t, derivation, "named", "id", "name") {
				if len(derivation.Constraints) != 1 {
					t.Logf("expected 1 constraint, found %d", len(derivation.Constraints))
					t.Fail()
				} else {
					assertIsTriplet(t, derivation.Constraints[0], Triplet{
						Id:    VarOf[Number]("id"),
						Attr:  Part(Attribute("names")),
						Value: VarOf[Literal]("name"),
					})
				}
			}
		}
	}
}

func assertIsType[T Ast](t *testing.T, elm Ast) (T, bool) {
	t.Helper()
	typed, matched := elm.(T)
	if !matched {
		t.Logf("expected type %T, found %s", typed, elm)
		t.Fail()
	}
	return typed, matched
}

func assertIsDecl(t *testing.T, actual DerivationDecl, name string, accept ...Variable) bool {
	t.Helper()
	matches := true
	if name != actual.Name {
		t.Logf("expected name %q, found %q", name, actual.Name)
		t.Fail()
		matches = false
	}

	if len(accept) != len(actual.Accepting) {
		t.Logf("expected %d inputs, found %d", len(accept), len(actual.Accepting))
		t.Fail()
		matches = false
	}

	for i := range min(len(accept), len(actual.Accepting)) {
		if accept[i] != actual.Accepting[i] {
			t.Logf("expected input %d to be %q, found %q", i, accept[i], actual.Accepting[i])
			t.Fail()
		}
	}

	return matches
}

func assertIsTriplet(t *testing.T, elm Ast, expected Triplet) bool {
	t.Helper()
	actual, isTriplet := elm.(Triplet)
	if !isTriplet {
		t.Logf("expected %s but found %s", expected, elm)
		t.Fail()
		return false
	}

	eq := actual.Eq(expected)

	if !eq {
		t.Logf("expected triplet %s, found %s", expected, actual)
		t.Fail()
	}
	return eq
}

func assertIsInvoke(t *testing.T, elm Ast, expected DerivationInvoke) bool {
	t.Helper()
	actual, isInvoke := elm.(DerivationInvoke)
	if !isInvoke {
		t.Logf("expected %s but found %s", expected, elm)
		t.Fail()
		return false
	}

	eq := actual.Eq(expected)

	if !eq {
		t.Logf("expected invoke %s, found %s", expected, actual)
		t.Fail()
	}
	return eq
}
