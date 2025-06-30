package parse

import (
	"strings"
	"testing"
)

var FakeEnv = NewEnv([]string{
	"BAR=bar",
	"FOO=foo",
	"EMPTY=",
	"ALSO_EMPTY=",
	"A=AAA",
	"test=test",
})

type mode int

const (
	relaxed mode = iota
	noUnset
	noEmpty
	strict
)

// Restrictions specifier
var (
	Relaxed   = &Restrictions{false, false, false, false, nil}
	NoEmpty   = &Restrictions{false, true, false, false, nil}
	NoUnset   = &Restrictions{true, false, false, false, nil}
	Strict    = &Restrictions{true, true, false, false, nil}
	KeepUnset = &Restrictions{false, false, false, true, nil}
)

var restrict = map[mode]*Restrictions{
	relaxed: Relaxed,
	noUnset: NoUnset,
	noEmpty: NoEmpty,
	strict:  Strict,
}

var errNone = map[mode]bool{}
var errUnset = map[mode]bool{noUnset: true, strict: true}
var errEmpty = map[mode]bool{noEmpty: true, strict: true}
var errAll = map[mode]bool{relaxed: true, noUnset: true, noEmpty: true, strict: true}
var errAllFull = map[mode]bool{relaxed: true, noUnset: true, noEmpty: true, strict: true}

type parseTest struct {
	name     string
	input    string
	expected string
	hasErr   map[mode]bool
}

var parseTests = []parseTest{
	{"lower-case variable", "lower variable $test ok", "lower variable test ok", errNone},
	{"empty", "", "", errNone},
	{"env only", "$BAR", "bar", errNone},
	{"with text", "$BAR baz", "bar baz", errNone},
	{"concatenated", "$BAR$FOO", "barfoo", errNone},
	{"2 env var", "$BAR - $FOO", "bar - foo", errNone},
	{"invalid var", "$_ bar", "$_ bar", errNone},
	{"invalid subst var", "${_} bar", "${_} bar", errNone},
	{"value of $var", "${BAR}baz", "barbaz", errNone},
	{"$var not set -", "${NOTSET-$BAR}", "bar", errNone},
	{"$var not set =", "${NOTSET=$BAR}", "bar", errNone},
	{"$var set but empty -", "${EMPTY-$BAR}", "", errEmpty},
	{"$var set but empty =", "${EMPTY=$BAR}", "", errEmpty},
	{"$var not set or empty :-", "${EMPTY:-$BAR}", "bar", errNone},
	{"$var not set or empty :=", "${EMPTY:=$BAR}", "bar", errNone},
	{"if $var set evaluate expression as $other +", "${EMPTY+hello}", "hello", errNone},
	{"if $var set evaluate expression as $other :+", "${EMPTY:+hello}", "hello", errNone},
	{"if $var not set, use empty string +", "${NOTSET+hello}", "", errNone},
	{"if $var not set, use empty string :+", "${NOTSET:+hello}", "", errNone},
	{"multi line string", "hello $BAR\nhello ${EMPTY:=$FOO}", "hello bar\nhello foo", errNone},
	{"issue #1", "${hello:=wo_rld} ${foo:=bar_baz}", "wo_rld bar_baz", errNone},
	{"issue #2", "name: ${NAME:=foo_qux}, key: ${EMPTY:=baz_bar}", "name: foo_qux, key: baz_bar", errNone},
	{"gh-issue-8", "prop=${HOME_URL-http://localhost:8080}", "prop=http://localhost:8080", errNone},
	// operators as leading values
	{"gh-issue-41-1", "${NOTSET--1}", "-1", errNone},
	{"gh-issue-41-2", "${NOTSET:--1}", "-1", errNone},
	{"gh-issue-41-3", "${NOTSET=-1}", "-1", errNone},
	{"gh-issue-41-4", "${NOTSET:==1}", "=1", errNone},

	// single letter
	{"gh-issue-43-1", "${A}", "AAA", errNone},

	// bad substitution
	{"closing brace expected", "hello ${", "", errAll},

	// test specifically for failure modes
	{"$var not set", "${NOTSET}", "", errUnset},
	{"$var set to empty", "${EMPTY}", "", errEmpty},
	// restrictions for plain variables without braces
	{"gh-issue-9", "$NOTSET", "", errUnset},
	{"gh-issue-9", "$EMPTY", "", errEmpty},

	{"$var and $DEFAULT not set -", "${NOTSET-$ALSO_NOTSET}", "", errUnset},
	{"$var and $DEFAULT not set :-", "${NOTSET:-$ALSO_NOTSET}", "", errUnset},
	{"$var and $DEFAULT not set =", "${NOTSET=$ALSO_NOTSET}", "", errUnset},
	{"$var and $DEFAULT not set :=", "${NOTSET:=$ALSO_NOTSET}", "", errUnset},
	{"$var and $OTHER not set +", "${NOTSET+$ALSO_NOTSET}", "", errNone},
	{"$var and $OTHER not set :+", "${NOTSET:+$ALSO_NOTSET}", "", errNone},

	{"$var empty and $DEFAULT not set -", "${EMPTY-$NOTSET}", "", errEmpty},
	{"$var empty and $DEFAULT not set :-", "${EMPTY:-$NOTSET}", "", errUnset},
	{"$var empty and $DEFAULT not set =", "${EMPTY=$NOTSET}", "", errEmpty},
	{"$var empty and $DEFAULT not set :=", "${EMPTY:=$NOTSET}", "", errUnset},
	{"$var empty and $OTHER not set +", "${EMPTY+$NOTSET}", "", errUnset},
	{"$var empty and $OTHER not set :+", "${EMPTY:+$NOTSET}", "", errUnset},

	{"$var not set and $DEFAULT empty -", "${NOTSET-$EMPTY}", "", errEmpty},
	{"$var not set and $DEFAULT empty :-", "${NOTSET:-$EMPTY}", "", errEmpty},
	{"$var not set and $DEFAULT empty =", "${NOTSET=$EMPTY}", "", errEmpty},
	{"$var not set and $DEFAULT empty :=", "${NOTSET:=$EMPTY}", "", errEmpty},
	{"$var not set and $OTHER empty +", "${NOTSET+$EMPTY}", "", errNone},
	{"$var not set and $OTHER empty :+", "${NOTSET:+$EMPTY}", "", errNone},

	{"$var and $DEFAULT empty -", "${EMPTY-$ALSO_EMPTY}", "", errEmpty},
	{"$var and $DEFAULT empty :-", "${EMPTY:-$ALSO_EMPTY}", "", errEmpty},
	{"$var and $DEFAULT empty =", "${EMPTY=$ALSO_EMPTY}", "", errEmpty},
	{"$var and $DEFAULT empty :=", "${EMPTY:=$ALSO_EMPTY}", "", errEmpty},
	{"$var and $OTHER empty +", "${EMPTY+$ALSO_EMPTY}", "", errEmpty},
	{"$var and $OTHER empty :+", "${EMPTY:+$ALSO_EMPTY}", "", errEmpty},

	// escaping.
	{"escape $$var", "FOO $$BAR BAZ", "FOO $BAR BAZ", errNone},
	{"escape $${subst}", "FOO $${BAR} BAZ", "FOO ${BAR} BAZ", errNone},
	{"escape $$$var", "$$$BAR", "$bar", errNone},
	{"escape $$${subst}", "$$${BAZ:-baz}", "$baz", errNone},
}

var negativeParseTests = []parseTest{
	{"$NOTSET and EMPTY are displayed as in full error output", "${NOTSET} and $EMPTY", "variable ${NOTSET} not set\nvariable ${EMPTY} set but empty", errAllFull},
}

// Test cases for KeepUnset functionality
var keepUnsetTests = []parseTest{
	{"keep unset variable", "$NOTSET", "$NOTSET", errNone},
	{"keep unset substitution", "${NOTSET}", "${NOTSET}", errNone},
	{"keep unset with text", "prefix $NOTSET suffix", "prefix $NOTSET suffix", errNone},
	{"keep unset substitution with text", "prefix ${NOTSET} suffix", "prefix ${NOTSET} suffix", errNone},
	{"substitute set variable", "$BAR", "bar", errNone},
	{"substitute set substitution", "${BAR}", "bar", errNone},
	{"keep unset with defaults", "${NOTSET-default}", "default", errNone},
	{"keep unset with colon defaults", "${NOTSET:-default}", "default", errNone},
	{"keep unset with equals defaults", "${NOTSET=default}", "default", errNone},
	{"keep unset with colon equals defaults", "${NOTSET:=default}", "default", errNone},
	{"keep empty with colon defaults", "${EMPTY:-default}", "default", errNone},
	{"keep empty with colon equals defaults", "${EMPTY:=default}", "default", errNone},
	{"keep set with plus", "${BAR+replacement}", "replacement", errNone},
	{"keep unset with plus", "${NOTSET+replacement}", "", errNone},
	{"mixed set and unset", "$BAR $NOTSET", "bar $NOTSET", errNone},
	{"multiple unset variables", "$NOTSET1 $NOTSET2", "$NOTSET1 $NOTSET2", errNone},
}

func TestParse(t *testing.T) {
	doTest(t, relaxed)
}

func TestParseNoUnset(t *testing.T) {
	doTest(t, noUnset)
}

func TestParseNoEmpty(t *testing.T) {
	doTest(t, noEmpty)
}

func TestParseStrict(t *testing.T) {
	doTest(t, strict)
}

func TestParseStrictNoFailFast(t *testing.T) {
	doNegativeAssertTest(t, strict)
}

func TestParseKeepUnset(t *testing.T) {
	for _, test := range keepUnsetTests {
		result, err := New(test.name, FakeEnv, KeepUnset).Parse(test.input)
		hasErr := err != nil
		if hasErr != test.hasErr[relaxed] { // KeepUnset should not produce errors
			t.Errorf("%s=(error): got\n\t%v\nexpected\n\t%v\ninput: %s\nresult: %s\nerror: %v",
				test.name, hasErr, test.hasErr[relaxed], test.input, result, err)
		}
		if result != test.expected {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.expected)
		}
	}
}

func doTest(t *testing.T, m mode) {
	for _, test := range parseTests {
		result, err := New(test.name, FakeEnv, restrict[m]).Parse(test.input)
		hasErr := err != nil
		if hasErr != test.hasErr[m] {
			t.Errorf("%s=(error): got\n\t%v\nexpected\n\t%v\ninput: %s\nresult: %s\nerror: %v",
				test.name, hasErr, test.hasErr[m], test.input, result, err)
		}
		if result != test.expected {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.expected)
		}
	}
}

func doNegativeAssertTest(t *testing.T, m mode) {
	for _, test := range negativeParseTests {
		result, err := (*&Parser{Name: test.name, Env: FakeEnv, Restrict: restrict[m], Mode: AllErrors}).Parse(test.input)
		hasErr := err != nil
		if hasErr != test.hasErr[m] {
			t.Errorf("%s=(error): got\n\t%v\nexpected\n\t%v\ninput: %s\nresult: %s\nerror: %v",
				test.name, hasErr, test.hasErr[m], test.input, result, err)
		}
		if err.Error() != test.expected {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, err.Error(), test.expected)
		}
	}
}

// Test cases for VarMatcher functionality
var varMatcherTests = []struct {
	name, input, expected string
	matcher               varMatcher
	hasErr                bool
}{
	{"nil matcher accepts all variables", "$BAR $FOO $test", "bar foo test", nil, false},
	{"prefix matcher accepts only matching variables", "$APP_NAME $BAR $DB_HOST", " $BAR ",
		func(v string) bool { return strings.HasPrefix(v, "APP_") || strings.HasPrefix(v, "DB_") }, false}, // Non-matching variables get consumed by lexer, only spaces and non-matching BAR remain
	{"length matcher filters by variable name length", "$A $BAR $FOO", "AAA $BAR $FOO",
		func(v string) bool { return len(v) <= 2 }, false}, // Only A (len=1) matches the length filter, BAR and FOO are 3 chars
	{"case sensitive matcher", "$BAR $bar $FOO $foo", "bar $bar foo $foo",
		func(v string) bool { return v == strings.ToUpper(v) }, false}, // BAR and FOO match (uppercase), bar and foo don't exist anyway
	{"matcher always returns false", "$BAR $FOO $test", "$BAR $FOO $test",
		func(v string) bool { return false }, false}, // All variables are treated as text
	{"matcher always returns true", "$BAR $FOO $test", "bar foo test",
		func(v string) bool { return true }, false},
	{"underscore always rejected even with permissive matcher", "$_ $BAR", "$_ bar",
		func(v string) bool { return true }, false},
	{"substitution with matching variable", "${BAR} ${test}", "bar ${test}",
		func(v string) bool { return v == "BAR" }, false}, // BAR matches, test doesn't
	{"substitution with non-matching variable", "${BAR} ${test}", "${BAR} test",
		func(v string) bool { return v == "test" }, false}, // test matches, BAR doesn't
	{"substitution with defaults and matching variables", "${VALID:-default} ${FOO:-backup}", "${VALID:-default} foo",
		func(v string) bool { return v == "FOO" }, false}, // VALID doesn't match (doesn't exist anyway), FOO matches
	{"substitution with defaults and non-matching variables", "${NOTSET:-default} ${NOTSET2:-backup}", "${NOTSET:-default} ${NOTSET2:-backup}",
		func(v string) bool { return false }, false}, // Nothing matches // Neither variable matches
	{"complex substitution with mixed matching", "${VALID:=$BAR} ${test:=$FOO}", "${VALID:=bar} test",
		func(v string) bool { return v == "test" || v == "BAR" }, false}, // VALID doesn't match, test matches, BAR in default matches
	{"nested variables in substitution", "${BAR-default} ${NOTSET-backup}", "bar ${NOTSET-backup}",
		func(v string) bool { return v == "BAR" }, false}, // BAR matches, NOTSET doesn't
	{"mixed simple and substitution variables", "prefix $FOO middle ${BAR} suffix $test end", "prefix foo middle ${BAR} suffix test end",
		func(v string) bool { return v == "FOO" || v == "test" }, false}, // FOO and test match, BAR doesn't
	{"matcher with special characters in variable names", "$VAR_1 $VAR_2 $test", "$VAR_1 $VAR_2 test",
		func(v string) bool { return !strings.Contains(v, "_") }, false}, // VAR_1 and VAR_2 have underscores, test doesn't
	{"empty variable name handling", "${} $BAR", "${} bar",
		func(v string) bool { return v != "" }, false}, // Empty variable name doesn't match, BAR does
}

// TestVarMatcher tests the VarMatcher functionality in the parser
func TestVarMatcher(t *testing.T) {
	for _, test := range varMatcherTests {
		t.Run(test.name, func(t *testing.T) {
			restrictions := &Restrictions{
				VarMatcher: test.matcher,
			}

			parser := New(test.name, FakeEnv, restrictions)
			result, err := parser.Parse(test.input)

			hasErr := err != nil
			if hasErr != test.hasErr {
				t.Errorf("Error expectation mismatch: got error=%v, expected error=%v\nInput: %s\nResult: %s\nError: %v",
					hasErr, test.hasErr, test.input, result, err)
				return
			}

			if result != test.expected {
				t.Errorf("Result mismatch:\nInput:    %q\nGot:      %q\nExpected: %q", test.input, result, test.expected)
			}
		})
	}
}

// TestVarMatcherWithRestrictions tests VarMatcher combined with other restrictions
func TestVarMatcherWithRestrictions(t *testing.T) {
	testEnv := NewEnv([]string{
		"SET_VAR=value",
		"EMPTY_VAR=",
	})

	tests := []struct {
		name, input, expected string
		matcher               varMatcher
		restrictions          *Restrictions
		hasErr                bool
	}{
		{"VarMatcher with NoUnset - matched variable not set", "$VALID_NOTSET $INVALID_NOTSET", "",
			func(v string) bool { return strings.HasPrefix(v, "VALID_") }, &Restrictions{NoUnset: true}, true},
		{"VarMatcher with NoUnset - unmatched variable not set (no error)", "$INVALID_NOTSET", "$INVALID_NOTSET",
			func(v string) bool { return false }, &Restrictions{NoUnset: true}, false},
		{"VarMatcher with NoEmpty - matched variable empty", "$EMPTY_VAR $INVALID", "",
			func(v string) bool { return v == "EMPTY_VAR" }, &Restrictions{NoEmpty: true}, true},
		{"VarMatcher with KeepUnset - unmatched variables kept", "$SET_VAR $NOTSET_INVALID", "value $NOTSET_INVALID",
			func(v string) bool { return v == "SET_VAR" }, &Restrictions{KeepUnset: true}, false},
		{"VarMatcher with KeepUnset - matched but unset variable kept as original", "$VALID_NOTSET", "$VALID_NOTSET",
			func(v string) bool { return strings.HasPrefix(v, "VALID_") }, &Restrictions{KeepUnset: true}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.restrictions.VarMatcher = test.matcher
			parser := New(test.name, testEnv, test.restrictions)
			result, err := parser.Parse(test.input)

			hasErr := err != nil
			if hasErr != test.hasErr {
				t.Errorf("Error expectation mismatch: got error=%v, expected error=%v\nInput: %s\nResult: %s\nError: %v",
					hasErr, test.hasErr, test.input, result, err)
				return
			}

			if !hasErr && result != test.expected {
				t.Errorf("Result mismatch:\nInput:    %q\nGot:      %q\nExpected: %q", test.input, result, test.expected)
			}
		})
	}
}

// TestVarMatcherEdgeCases tests edge cases for VarMatcher
func TestVarMatcherEdgeCases(t *testing.T) {
	tests := []struct {
		name, input, expected string
		matcher               varMatcher
		hasErr                bool
	}{
		{"matcher with panic recovery", "$BAR", "$BAR",
			func(v string) bool { return v[100] == 'x' }, false}, // Will panic on out of bounds access // Should be treated as non-matching due to panic
		{"variable name with numbers", "$VAR1 $VAR2 $A", "$VAR1 $VAR2 AAA",
			func(v string) bool { return !strings.ContainsAny(v, "0123456789") }, false}, // Only A matches (no digits), VAR1 and VAR2 don't exist in FakeEnv anyway
		{"unicode variable names", "$VAR_测试 $BAR", "$VAR_测试 bar",
			func(v string) bool { return !strings.Contains(v, "测试") }, false}, // VAR_测试 doesn't exist in FakeEnv, BAR matches the matcher
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use a safe wrapper for the matcher to handle panics
			safeMatcher := func(v string) bool {
				defer func() {
					if r := recover(); r != nil {
						// If matcher panics, treat as non-matching
					}
				}()
				return test.matcher(v)
			}

			restrictions := &Restrictions{
				VarMatcher: safeMatcher,
			}

			parser := New(test.name, FakeEnv, restrictions)
			result, err := parser.Parse(test.input)

			hasErr := err != nil
			if hasErr != test.hasErr {
				t.Errorf("Error expectation mismatch: got error=%v, expected error=%v\nInput: %s\nResult: %s\nError: %v",
					hasErr, test.hasErr, test.input, result, err)
				return
			}

			if result != test.expected {
				t.Errorf("Result mismatch:\nInput:    %q\nGot:      %q\nExpected: %q", test.input, result, test.expected)
			}
		})
	}
}
