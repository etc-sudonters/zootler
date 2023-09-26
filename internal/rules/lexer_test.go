package rules

import (
	"testing"
	"unicode/utf8"

	"github.com/etc-sudonters/zootler/internal/testutils"
)

func TestCanLexIdentifier(t *testing.T) {
	rule := "Compiler"
	expected := []item{{typ: itemIdent, val: rule}}

	l := lex("TestCanLexIdentifier", rule)
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsIfNonSepFollowsIdent(t *testing.T) {
	expected := []item{
		{typ: itemErr, pos: 8, val: "unexpected '\\''"},
	}

	l := lex("TestErrsIfNonSepFollowsIdent", "Compiler'")
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexBoolOp(t *testing.T) {
	rule := "Compiler or Interpreter"
	l := lex("TestCanLexBoolOp", rule)
	expected := []item{
		{typ: itemIdent, val: "Compiler", pos: 0},
		{typ: itemBoolOp, val: "or", pos: 9},
		{typ: itemIdent, val: "Interpreter", pos: 12},
	}

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexNumber(t *testing.T) {
	l := lex("TestCanLexNumber", "99")
	expected := []item{{typ: itemNumber, pos: 0, val: "99"}}

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsIfNonSepFollowsNumber(t *testing.T) {
	expected := []item{
		{typ: itemErr, pos: 2, val: "unexpected '\\''"},
	}

	l := lex("TestErrsIfNonSepFollowsNumber", "99'")

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexParens(t *testing.T) {
	rule := "(Compiler or Interpreter)"
	expected := []item{
		{typ: itemOpenParen, pos: 0, val: "("},
		{typ: itemIdent, pos: 1, val: "Compiler"},
		{typ: itemBoolOp, pos: 10, val: "or"},
		{typ: itemIdent, pos: 13, val: "Interpreter"},
		{typ: itemCloseParen, pos: 24, val: ")"},
	}

	l := lex("TestCanLexParens", rule)
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsOnUnclosedParen(t *testing.T) {
	rule := "(Compiler"
	expected := []item{
		{typ: itemOpenParen, pos: 0, val: "("},
		{typ: itemIdent, pos: 1, val: "Compiler"},
		{typ: itemErr, pos: 9, val: "unclosed '('"},
	}

	l := lex("TestErrsOnUnclosedParen", rule)
	l.debug = t.Logf
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsOnUnexpectedClosedParent(t *testing.T) {
	expected := []item{
		{typ: itemIdent, pos: 0, val: "Compiler"},
		{typ: itemErr, pos: 9, val: "unexpected ')'"},
	}

	l := lex("TestErrsOnUnexpectedClosedParent", "Compiler)")
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexActualRules(t *testing.T) {
	rules := []string{
		"can_play(Song_of_Time) or (logic_shadow_mq_invisible_blades and damage_multiplier != 'ohko')",
		"Boomerang or (logic_spirit_fence_gs and Kokiri_Sword) or (((Silver_Rupee_Spirit_Temple_Child_Early_Torches, 5) or logic_spirit_fence_gs) and (Sticks or has_explosives or Slingshot or can_use(Dins_Fire)))",
		"((can_use(Megaton_Hammer) and logic_dc_hammer_floor) or has_explosives or king_dodongo_shortcuts) and (((Bombs or Progressive_Strength_Upgrade) and can_jumpslash) or deadly_bonks == 'ohko')",
		"is_child and at_day and (can_break_crate or chicken_count < 7)",
	}

	for _, rule := range rules {
		l := lex("TestLexActualRule", rule)
		collected := lexUntilEofOrErr(l, t)

		if collected[len(collected)-1].typ == itemErr {
			t.Logf("Failed lexing rule %.80q...", rule)
			t.Log("Failure", collected[len(collected)-1])
			t.Log("Collected tokens", collected)
			t.Logf("lexer state %+v", l)
			t.Fail()
		}
	}

}

func lexUntilEofOrErr(l *lexer, t *testing.T) []item {
	// protection against lexer getting stuck
	runeCount := utf8.RuneCountInString(l.input)
	i := 0
	collected := []item{}

	for {
		if i > runeCount {
			t.Logf("lexer %s spinning, killing test", l.name)
			t.Logf("collected tokens %+v", collected)
			t.Logf("lexer state %+v", l)
			t.FailNow()
		}

		item := l.nextItem()
		if item.typ == itemEof {
			break
		}

		collected = append(collected, item)

		// collect err but not eof
		if item.typ == itemErr {
			break
		}

		i++
	}

	return collected
}

func toksAreEqual(expected, actual []item, t *testing.T) {
	testutils.ArrEqF(expected, actual, func(e, a item) bool { return e.fungible(a) }, t)
}

func (i item) fungible(o item) bool {
	return i.typ == o.typ && i.pos == o.pos && i.val == o.val
}
