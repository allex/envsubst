package parse

import (
	"fmt"
	"strings"
	"testing"
)

// TestPatternTransformerAbstraction demonstrates how the pattern transformer
// abstraction makes it easy to add new transformation patterns and verifies
// both normal transformation and KeepUnset behavior
func TestPatternTransformerAbstraction(t *testing.T) {
	testCases := []struct {
		name, input, expected string
		keepUnset             bool
		envVars               []string
	}{
		// Normal transformation tests
		{"uppercase transformation", "${TEST_VAR^^}", "HELLO WORLD", false, []string{"TEST_VAR=Hello World"}},
		{"lowercase transformation", "${TEST_VAR,,}", "hello world", false, []string{"TEST_VAR=Hello World"}},

		// KeepUnset behavior tests
		{"keep unset uppercase", "${UNSET_VAR^^}", "${UNSET_VAR^^}", true, []string{"SET_VAR=Hello"}},
		{"keep unset lowercase", "${UNSET_VAR,,}", "${UNSET_VAR,,}", true, []string{"SET_VAR=Hello"}},
		{"transform set uppercase with keepUnset", "${SET_VAR^^}", "HELLO", true, []string{"SET_VAR=Hello"}},
		{"transform set lowercase with keepUnset", "${SET_VAR,,}", "hello", true, []string{"SET_VAR=Hello"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := NewEnv(tc.envVars)
			parser := New("test", env, &Restrictions{KeepUnset: tc.keepUnset})

			result, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

// TestPatternTransformerRegistry demonstrates the RegisterPatternTransformer function
func TestPatternTransformerRegistry(t *testing.T) {
	// Save original state
	originalDefinitions := make(map[itemType]PatternDefinition)
	for k, v := range patternDefinitions {
		originalDefinitions[k] = v
	}

	// Restore original state after test
	defer func() {
		patternDefinitions = originalDefinitions
	}()

	// Create a hypothetical new itemType for demonstration
	// (In real usage, this would be defined in lex.go)
	const itemTitleCase itemType = 999

	// Register a new pattern transformer for title case
	RegisterPatternTransformer(itemTitleCase, "~T", strings.Title)

	// Verify it was registered
	if patternDef, exists := patternDefinitions[itemTitleCase]; !exists {
		t.Error("expected pattern definition to be registered")
	} else {
		result := patternDef.Transformer("hello world")
		expected := "Hello World"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
		if patternDef.Operator != "~T" {
			t.Errorf("expected operator %q, got %q", "~T", patternDef.Operator)
		}
	}
}

// TestPatternDefinitionStructCohesion demonstrates the cohesiveness of PatternDefinition
func TestPatternDefinitionStructCohesion(t *testing.T) {
	// Test that PatternDefinition keeps related data together
	testCases := []struct {
		name     string
		itemType itemType
		input    string
		expected string
	}{
		{"uppercase pattern", itemCaretCaret, "hello world", "HELLO WORLD"},
		{"lowercase pattern", itemCommaComma, "HELLO WORLD", "hello world"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patternDef, exists := patternDefinitions[tc.itemType]
			if !exists {
				t.Fatalf("pattern definition not found for %v", tc.itemType)
			}

			// Test transformer function
			result := patternDef.Transformer(tc.input)
			if result != tc.expected {
				t.Errorf("transformer: expected %q, got %q", tc.expected, result)
			}

			// Test suffix is properly defined
			if patternDef.Operator == "" {
				t.Error("operator should not be empty")
			}

			// Verify operator matches expected format for the pattern
			switch tc.itemType {
			case itemCaretCaret:
				if patternDef.Operator != "^^" {
					t.Errorf("expected operator ^^, got %q", patternDef.Operator)
				}
			case itemCommaComma:
				if patternDef.Operator != ",," {
					t.Errorf("expected operator ,,, got %q", patternDef.Operator)
				}
			}

			// Test that both components work together for a complete pattern
			varName := "TEST_VAR"
			expectedPattern := "${" + varName + patternDef.Operator + "}"
			t.Logf("Complete pattern: %s", expectedPattern)
		})
	}
}

// TestPatternDefinitionCompleteness verifies all pattern definitions are complete
func TestPatternDefinitionCompleteness(t *testing.T) {
	for itemType, patternDef := range patternDefinitions {
		t.Run(fmt.Sprintf("itemType_%d", itemType), func(t *testing.T) {
			// Verify transformer is not nil
			if patternDef.Transformer == nil {
				t.Error("transformer function should not be nil")
			}

			// Verify suffix is not empty
			if patternDef.Operator == "" {
				t.Error("operator should not be empty")
			}

			// Test transformer works
			testInput := "Test Input"
			result := patternDef.Transformer(testInput)
			if result == "" && testInput != "" {
				t.Error("transformer should not return empty string for non-empty input")
			}
		})
	}
}
