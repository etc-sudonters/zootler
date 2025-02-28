package query

import (
	"sudonters/libzootr/galoshes/parse"
	"testing"

	"github.com/etc-sudonters/substrate/peruse"
)

func TestCanParseAndPlanFind(t *testing.T) {
	findScript := `find [ $place-id $place-name ]
where [
    [ :holding $place-id $token-id ]
    [ :is-named $place-id $place-name ]
    [ :is-named $token-id "Kokri Sword" ]
    [ $token-id tokens 1 ] 
]
rules [
    [:holding $place $token] [ [$place world/placement/holds $token] ]
    [:named $id $name] [ [$id global/name $name] ]
    [:is-named $id $name] [ [:has-name $id $name]]
    [:has-name $id $name] [ [:named $id $name]]
]
`

	lexer := parse.NewLexer(findScript)
	parser := peruse.NewParser(parse.NewGrammar(), lexer)

	result, parseErr := parser.ParseAt(peruse.LOWEST)
	if parseErr != nil {
		t.Log(findScript)
		t.Logf("lexer state: %#v", lexer)
		t.Logf("parser state: %#v", parser)
		t.Fatal(parseErr.Error())
	}

	find, isFind := result.(*parse.FindNode)

	if !isFind {
		t.Log(findScript)
		t.Logf("%#v", result)
		t.Fatal("failed to parse find")
	}

	ta := parse.NewAnnotator()
	ta.VisitFindNode(find)
	sub := parse.NewSubstituter(ta.Substitutions)
	sub.VisitFindNode(find)

	plan := BuildQueryPlan(find)
	t.Logf("%v", plan)
	t.Fail()
}
