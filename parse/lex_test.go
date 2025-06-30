package parse

import (
	"strings"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

type lexTestWithMatcher struct {
	name    string
	input   string
	matcher varMatcher
	want    []item
}

var (
	tEOF       = item{itemEOF, 0, ""}
	tPlus      = item{itemPlus, 0, ""}
	tDash      = item{itemDash, 0, "-"}
	tEquals    = item{itemEquals, 0, "="}
	tColEquals = item{itemColonEquals, 0, ":="}
	tColDash   = item{itemColonDash, 0, ":-"}
	tColPlus   = item{itemColonPlus, 0, ":+"}
	tLeft      = item{itemLeftDelim, 0, "${"}
	tRight     = item{itemRightDelim, 0, "}"}
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"text", "hello", []item{
		{itemText, 0, "hello"},
		tEOF,
	}},
	{"var", "$hello", []item{
		{itemVariable, 0, "$hello"},
		tEOF,
	}},
	{"single char var", "${A}", []item{
		tLeft,
		{itemVariable, 0, "A"},
		tRight,
		tEOF,
	}},
	{"2 vars", "$hello $world", []item{
		{itemVariable, 0, "$hello"},
		{itemText, 0, " "},
		{itemVariable, 0, "$world"},
		tEOF,
	}},
	{"substitution-1", "bar ${BAR}", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tRight,
		tEOF,
	}},
	{"substitution-2", "bar ${BAR:=baz}", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tColEquals,
		{itemText, 0, "b"},
		{itemText, 0, "a"},
		{itemText, 0, "z"},
		tRight,
		tEOF,
	}},
	{"substitution-3", "bar ${BAR:=$BAZ}", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tColEquals,
		{itemVariable, 0, "$BAZ"},
		tRight,
		tEOF,
	}},
	{"substitution-4", "bar ${BAR:=$BAZ} foo", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tColEquals,
		{itemVariable, 0, "$BAZ"},
		tRight,
		{itemText, 0, " foo"},
		tEOF,
	}},
	{"substitution-leading-dash-1", "bar ${BAR:--1} foo", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tColDash,
		{itemText, 0, "-"},
		{itemText, 0, "1"},
		tRight,
		{itemText, 0, " foo"},
		tEOF,
	}},
	{"substitution-leading-dash-2", "bar ${BAR:=-1} foo", []item{
		{itemText, 0, "bar "},
		tLeft,
		{itemVariable, 0, "BAR"},
		tColEquals,
		{itemText, 0, "-"},
		{itemText, 0, "1"},
		tRight,
		{itemText, 0, " foo"},
		tEOF,
	}},
	{"closing brace error", "hello-${world", []item{
		{itemText, 0, "hello-"},
		tLeft,
		{itemVariable, 0, "world"},
		{itemError, 0, "closing brace expected"},
	}},
	{"escaping $$var", "hello $$HOME", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "$"},
		{itemText, 8, "HOME"},
		tEOF,
	}},
	{"escaping $${subst}", "hello $${HOME}", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "$"},
		{itemText, 8, "{HOME}"},
		tEOF,
	}},
	{"no digit $1", "hello $1", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "$1"},
		tEOF,
	}},
	{"no digit $1ABC", "hello $1ABC", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "$1"},
		{itemText, 9, "ABC"},
		tEOF,
	}},
	{"no digit ${2}", "hello ${2}", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "${2"},
		{itemText, 10, "}"},
		tEOF,
	}},
	{"no digit ${2ABC}", "hello ${2ABC}", []item{
		{itemText, 0, "hello "},
		{itemText, 7, "${2"},
		{itemText, 10, "ABC}"},
		tEOF,
	}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s:\ninput\n\t%q\ngot\n\t%+v\nexpected\n\t%v", test.name, test.input, items, test.items)
		}
	}
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	noDigit := strings.HasPrefix(t.name, "no digit")
	l := lex(t.input, noDigit, nil)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

// collectWithMatcher gathers the emitted items into a slice using a custom matcher.
func collectWithMatcher(t *lexTest, matcher varMatcher) (items []item) {
	noDigit := strings.HasPrefix(t.name, "no digit")
	l := lex(t.input, noDigit, matcher)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

// TestLexMatcherPatterns tests various matcher patterns
func TestLexMatcherPatterns(t *testing.T) {
	tests := []lexTestWithMatcher{
		{"nil matcher accepts all", "$hello $world", nil, []item{
			{itemVariable, 0, "$hello"},
			{itemText, 0, " "},
			{itemVariable, 0, "$world"},
			tEOF,
		}},
		{"prefix matcher", "$APP_NAME $DB_HOST $USER_ID", func(v string) bool {
			return strings.HasPrefix(v, "APP_") || strings.HasPrefix(v, "DB_")
		}, []item{
			{itemVariable, 0, "$APP_NAME"},
			{itemText, 0, " "},
			{itemVariable, 0, "$DB_HOST"},
			{itemText, 0, " "},
			{itemText, 0, "$USER_ID"},
			tEOF,
		}},
		{"length matcher", "$A $AB $ABC $ABCD", func(v string) bool {
			return len(v) <= 2
		}, []item{
			{itemVariable, 0, "$A"},
			{itemText, 0, " "},
			{itemVariable, 0, "$AB"},
			{itemText, 0, " "},
			{itemText, 0, "$ABC"},
			{itemText, 0, " "},
			{itemText, 0, "$ABCD"},
			tEOF,
		}},
		{"case sensitive matcher", "$hello $HELLO $Hello", func(v string) bool {
			return v == strings.ToUpper(v)
		}, []item{
			{itemText, 0, "$hello"},
			{itemText, 0, " "},
			{itemVariable, 0, "$HELLO"},
			{itemText, 0, " "},
			{itemText, 0, "$Hello"},
			tEOF,
		}},
		{"substitution with matcher", "${VALID} ${invalid}", func(v string) bool {
			return v == strings.ToUpper(v)
		}, []item{
			tLeft,
			{itemVariable, 0, "VALID"},
			tRight,
			{itemText, 0, " "},
			tLeft,
			{itemText, 0, "invalid"},
			tRight,
			tEOF,
		}},
		{"mixed vars", "$hello $world $foo", func(v string) bool {
			return v == "world"
		}, []item{
			{itemText, 0, "$hello"},
			{itemText, 0, " "},
			{itemVariable, 0, "$world"},
			{itemText, 0, " "},
			{itemText, 0, "$foo"},
			tEOF,
		}},
		{"complex substitution with mixed matching", "${VALID:=$BAR} ${test:=$FOO}", func(v string) bool {
			return v == "test" || v == "BAR"
		}, []item{
			tLeft,
			{itemText, 0, "VALID"},
			tColEquals,
			{itemVariable, 0, "$BAR"},
			tRight,
			{itemText, 0, " "},
			tLeft,
			{itemVariable, 0, "test"},
			tColEquals,
			{itemText, 0, "$FOO"},
			tRight,
			tEOF,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := lexTest{name: tt.name, input: tt.input, items: tt.want}
			items := collectWithMatcher(&test, tt.matcher)
			if !equal(items, tt.want, false) {
				t.Errorf("TestLexMatcherPatterns %s:\ninput\n\t%q\ngot\n\t%+v\nexpected\n\t%v", tt.name, tt.input, items, tt.want)
			}
		})
	}
}

// TestLexMatcherEdgeCases tests edge cases for the matcher functionality
func TestLexMatcherEdgeCases(t *testing.T) {
	tests := []lexTestWithMatcher{
		{"empty variable name", "${}", func(v string) bool {
			return true
		}, []item{
			tLeft,
			tRight,
			tEOF,
		}},
		{"nested substitution", "${VAR:=${INNER}}", func(v string) bool {
			return v == "VAR" || v == "INNER"
		}, []item{
			tLeft,
			{itemVariable, 0, "VAR"},
			tColEquals,
			{itemText, 0, "$"},
			{itemText, 0, "{INNER"},
			tRight,
			{itemText, 0, "}"},
			tEOF,
		}},
		{"substitution with accepted var", "${world:=$backup}", func(v string) bool {
			return strings.HasPrefix(v, "world")
		}, []item{
			tLeft,
			{itemVariable, 0, "world"},
			tColEquals,
			{itemText, 0, "$backup"},
			tRight,
			tEOF,
		}},
		{"matcher always false", "$hello $world ${test}", func(v string) bool {
			return false
		}, []item{
			{itemText, 0, "$hello"},
			{itemText, 0, " "},
			{itemText, 0, "$world"},
			{itemText, 0, " "},
			tLeft,
			{itemText, 0, "test"},
			tRight,
			tEOF,
		}},
		{"underscore always rejected", "$_", func(v string) bool {
			return true // even if matcher returns true, underscore is always rejected
		}, []item{
			{itemText, 0, "$_"},
			tEOF,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := lexTest{name: tt.name, input: tt.input, items: tt.want}
			items := collectWithMatcher(&test, tt.matcher)
			if !equal(items, tt.want, false) {
				t.Errorf("TestLexMatcherEdgeCases %s:\ninput\n\t%q\ngot\n\t%+v\nexpected\n\t%v", tt.name, tt.input, items, tt.want)
			}
		})
	}
}
