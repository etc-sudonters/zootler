package rulesparser

import "testing"

func TestParseRealRule(t *testing.T) {
	r := "can_play(Song_of_Time) or (logic_shadow_mq_invisible_blades and damage_multiplier != 'ohko')"
	l := NewRulesLexer(r)
	p := NewRulesParser(l)

	rule, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if rule == nil {
		t.Fatal("did not parse rule")
	}

	if p.HasMore() {
		t.Fatal("trailing unparsed content")
	}

	t.Fatal("task failed successfully")
}

func TestParseConstRule(t *testing.T) {
	inputs := []struct {
		raw      string
		expected bool
	}{
		{raw: "True", expected: true},
		{raw: "False", expected: false},
	}

	for _, i := range inputs {
		i := i
		t.Run("Const"+i.raw, func(t *testing.T) {
			l := NewRulesLexer(i.raw)
			p := NewRulesParser(l)

			rule, err := p.Parse()
			if err != nil {
				t.Fatalf("expected to successfully parse '%s': %s", i.raw, err)
			}

			switch r := rule.(type) {
			case *Boolean:
				if r.Value != i.expected {
					t.Logf("expected to parse %s to ConstRule{ %t }", i.raw, i.expected)
					t.Logf("instead parsed to %v", r)
					t.FailNow()
				}
				break
			default:
				t.Fatalf("expected to parse 'True' to ConstRule not %v", rule)
				break
			}
		})
	}
}
