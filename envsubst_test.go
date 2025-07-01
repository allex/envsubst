package envsubst

import (
	"os"
	"testing"

	"github.com/allex/envsubst/parse"
)

func init() {
	os.Setenv("BAR", "bar")
}

// Basic integration tests. because we  already test the
// templating processing in envsubst/parse;
func TestIntegration(t *testing.T) {
	input, expected := "foo $BAR", "foo bar"
	str, err := String(input)
	if str != expected || err != nil {
		t.Error("Expect string integration test to pass")
	}
	bytes, err := Bytes([]byte(input))
	if string(bytes) != expected || err != nil {
		t.Error("Expect bytes integration test to pass")
	}
	bytes, err = ReadFile("testdata/file.tmpl")
	fexpected, err := os.ReadFile("testdata/file.out")
	if string(bytes) != string(fexpected) || err != nil {
		t.Error("Expect ReadFile integration test to pass")
	}
}

func TestKeepUnsetIntegration(t *testing.T) {
	// Test that undefined variables are kept as original text
	input := "foo $UNDEFINED_VAR ${ALSO_UNDEFINED} $BAR"
	expected := "foo $UNDEFINED_VAR ${ALSO_UNDEFINED} bar"

	str, err := StringRestrictedKeepUnset(input, false, false, false, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if str != expected {
		t.Errorf("Expected %q, got %q", expected, str)
	}

	// Test bytes function
	bytes, err := BytesRestrictedKeepUnset([]byte(input), false, false, false, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(bytes) != expected {
		t.Errorf("Expected %q, got %q", expected, string(bytes))
	}
}

func TestEnvInitializeAndLazyInjection(t *testing.T) {
	testCases := []struct {
		name, input, expected string
		initialEnvs           []string
		lazyVars              map[string]string
	}{
		{
			"initial and lazy vars",
			"$INIT_VAR: $LAZY_VAR",
			"initial: lazy",
			[]string{"INIT_VAR=initial", "USER=testuser"},
			map[string]string{"LAZY_VAR": "lazy"},
		},
		{
			"multiple initial and lazy",
			"${HOME}/${USER}/$PROJECT/$TASK",
			"/home/user/testuser/myproject/mytask",
			[]string{"HOME=/home/user", "USER=testuser", "PROJECT=myproject"},
			map[string]string{"TASK": "mytask"},
		},
		{
			"lazy override initial",
			"Config: $CONFIG",
			"Config: custom",
			[]string{"CONFIG=default", "OTHER=value"},
			map[string]string{"CONFIG": "custom"},
		},
		{
			"demonstrate override behavior",
			"$VALUE is the value",
			"overridden is the value",
			[]string{"VALUE=original", "KEEP=unchanged"},
			map[string]string{"VALUE": "overridden"},
		},
		{
			"mixed substitution order",
			"$A-$B-$C-$D",
			"1-2-3-4",
			[]string{"A=1", "C=3"},
			map[string]string{"B": "2", "D": "4"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create env with initial variables
			env := parse.NewEnv(tc.initialEnvs)

			// Parse using the env with both initial and lazy variables
			parser := parse.New("test", env, &parse.Restrictions{})

			// Inject additional variables lazily
			for key, value := range tc.lazyVars {
				env.Set(key, value)
			}

			result, err := parser.Parse(tc.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}

			// Verify that both initial and lazy variables are accessible
			for key, expectedVal := range tc.lazyVars {
				if got := env.Get(key); got != expectedVal {
					t.Errorf("Lazy var %s: expected %q, got %q", key, expectedVal, got)
				}
			}
		})
	}
}
