package rules

import (
	"testing"
	"unicode/utf8"

	"sudonters/zootler/internal/testutils"
)

func TestCanLexIdentifier(t *testing.T) {
	rule := "Compiler"
	expected := []Item{{Type: ItemIdentifier, Value: rule}}

	l := NewLexer("TestCanLexIdentifier", rule)
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsIfNonSepFollowsIdent(t *testing.T) {
	expected := []Item{
		{Type: ItemErr, Pos: 8, Value: "unexpected '\\''"},
	}

	l := NewLexer("TestErrsIfNonSepFollowsIdent", "Compiler'")
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexBoolOp(t *testing.T) {
	rule := "Compiler or Interpreter"
	l := NewLexer("TestCanLexBoolOp", rule)
	expected := []Item{
		{Type: ItemIdentifier, Value: "Compiler", Pos: 0},
		{Type: ItemOr, Value: "or", Pos: 9},
		{Type: ItemIdentifier, Value: "Interpreter", Pos: 12},
	}

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexNumber(t *testing.T) {
	l := NewLexer("TestCanLexNumber", "99")
	expected := []Item{{Type: ItemNumber, Pos: 0, Value: "99"}}

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsIfNonSepFollowsNumber(t *testing.T) {
	expected := []Item{
		{Type: ItemErr, Pos: 2, Value: "unexpected '\\''"},
	}

	l := NewLexer("TestErrsIfNonSepFollowsNumber", "99'")

	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestCanLexParens(t *testing.T) {
	rule := "(Compiler or Interpreter)"
	expected := []Item{
		{Type: ItemOpenParen, Pos: 0, Value: "("},
		{Type: ItemIdentifier, Pos: 1, Value: "Compiler"},
		{Type: ItemOr, Pos: 10, Value: "or"},
		{Type: ItemIdentifier, Pos: 13, Value: "Interpreter"},
		{Type: ItemCloseParen, Pos: 24, Value: ")"},
	}

	l := NewLexer("TestCanLexParens", rule)
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsOnUnclosedParen(t *testing.T) {
	rule := "(Compiler"
	expected := []Item{
		{Type: ItemOpenParen, Pos: 0, Value: "("},
		{Type: ItemIdentifier, Pos: 1, Value: "Compiler"},
		{Type: ItemErr, Pos: 9, Value: "unclosed '(' or '['"},
	}

	l := NewLexer("TestErrsOnUnclosedParen", rule)
	collected := lexUntilEofOrErr(l, t)

	toksAreEqual(expected, collected, t)
}

func TestErrsOnUnexpectedClosedParent(t *testing.T) {
	expected := []Item{
		{Type: ItemIdentifier, Pos: 0, Value: "Compiler"},
		{Type: ItemErr, Pos: 9, Value: "unexpected ')'"},
	}

	l := NewLexer("TestErrsOnUnexpectedClosedParent", "Compiler)")
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
		l := NewLexer("TestLexActualRule", rule)
		collected := lexUntilEofOrErr(l, t)

		if collected[len(collected)-1].Type == ItemErr {
			t.Logf("Failed lexing rule %.80q...", rule)
			t.Log("Failure", collected[len(collected)-1])
			t.Log("Collected tokens", collected)
			t.Logf("lexer state %+v", l)
			t.Fail()
		}
	}

}

func lexUntilEofOrErr(l *lexer, t *testing.T) []Item {
	// protection against lexer getting stuck
	runeCount := utf8.RuneCountInString(l.input)
	i := 0
	collected := []Item{}

	for {
		if i > runeCount {
			t.Logf("lexer %s spinning, killing test", l.name)
			t.Logf("collected tokens %+v", collected)
			t.Logf("lexer state %+v", l)
			t.FailNow()
		}

		item := l.nextItem()
		if item.Type == ItemEof {
			break
		}

		collected = append(collected, item)

		// collect err but not eof
		if item.Type == ItemErr {
			break
		}

		i++
	}

	return collected
}

func toksAreEqual(expected, actual []Item, t *testing.T) {
	testutils.ArrEqF(expected, actual, func(e, a Item) bool { return e.fungible(a) }, t)
}

func (i Item) fungible(o Item) bool {
	return i.Type == o.Type && i.Pos == o.Pos && i.Value == o.Value
}
