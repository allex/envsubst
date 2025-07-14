package parse

import (
	"fmt"
	"strings"
)

// PatternTransformer defines a function that transforms a variable value according to a specific pattern
type PatternTransformer func(value string) string

// PatternDefinition combines a transformer function with its syntax suffix
type PatternDefinition struct {
	Operator    string             // Bash expansion operator syntax (e.g., "^^", ",,")
	Transformer PatternTransformer // Function to transform the variable value
}

// Pattern Transformer System
//
// The pattern transformer system provides a generic way to handle bash variable expansion patterns
// that transform variable values. This abstraction makes it easy to add new transformation patterns
// without modifying the core parsing logic.
//
// Architecture:
// - PatternTransformer: Function type that defines how to transform values
// - PatternDefinition: Struct combining transformer function and operator syntax
// - patternDefinitions: Maps itemType to PatternDefinition structs
// - RegisterPatternTransformer: Helper function to register new patterns
//
// Adding New Patterns:
// 1. Define a new itemType in lex.go (e.g., itemTitleCase)
// 2. Add lexer support for the pattern in lexSubstitutionOperator
// 3. Register the pattern using RegisterPatternTransformer
//
// Example:
//   RegisterPatternTransformer(itemTitleCase, "~T", strings.Title)
//
// This would enable ${VAR~T} to convert variables to title case.

// patternDefinitions maps itemType to their corresponding pattern definitions
var patternDefinitions = map[itemType]PatternDefinition{
	itemCaretCaret: {"^^", strings.ToUpper}, // ^^ converts to uppercase
	itemCommaComma: {",,", strings.ToLower}, // ,, converts to lowercase
}

// RegisterPatternTransformer allows registering new pattern transformers
// This makes it easy to extend the system with additional transformation patterns
func RegisterPatternTransformer(itemType itemType, operator string, transformer PatternTransformer) {
	patternDefinitions[itemType] = PatternDefinition{operator, transformer}
}

type Node interface {
	Type() NodeType
	String() (string, error)
}

// NodeType identifies the type of a node.
type NodeType int

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText NodeType = iota
	NodeSubstitution
	NodeVariable
)

type TextNode struct {
	NodeType
	Text string
}

func NewText(text string) *TextNode {
	return &TextNode{NodeText, text}
}

func (t *TextNode) String() (string, error) {
	return t.Text, nil
}

type VariableNode struct {
	NodeType
	Ident    string // Variable identifier name (e.g., "VAR" from "$VAR" or "${VAR}")
	Env      *Env
	Restrict *Restrictions
}

func NewVariable(ident string, env *Env, restrict *Restrictions) *VariableNode {
	return &VariableNode{NodeVariable, ident, env, restrict}
}

func (t *VariableNode) String() (string, error) {
	// If KeepUnset is enabled and variable is not set, return source text
	if t.Restrict.KeepUnset && !t.isSet() {
		// Construct the source text format from ident
		return "$" + t.Ident, nil
	}

	if err := t.validateNoUnset(); err != nil {
		return "", err
	}
	value := t.Env.Get(t.Ident)
	if err := t.validateNoEmpty(value); err != nil {
		return "", err
	}
	return value, nil
}

func (t *VariableNode) isSet() bool {
	return t.Env.Has(t.Ident)
}

func (t *VariableNode) validateNoUnset() error {
	if t.Restrict.NoUnset && !t.isSet() {
		return Error(fmt.Sprintf("variable ${%s} not set", t.Ident), "NoUnset")
	}
	return nil
}

func (t *VariableNode) validateNoEmpty(value string) error {
	if t.Restrict.NoEmpty && value == "" && t.isSet() {
		return Error(fmt.Sprintf("variable ${%s} set but empty", t.Ident), "NoEmpty")
	}
	return nil
}

type SubstitutionNode struct {
	NodeType
	ExpType  itemType
	Variable *VariableNode
	Default  Node // Default could be variable or text
}

func (t *SubstitutionNode) String() (string, error) {
	// Handle pattern transformations using the transformer map
	if patternDef, hasPatternDef := patternDefinitions[t.ExpType]; hasPatternDef {
		if t.Variable.Restrict.KeepUnset && !t.Variable.isSet() {
			// Return original syntax for unset variables when KeepUnset is enabled
			return "${" + t.Variable.Ident + patternDef.Operator + "}", nil
		}

		value, err := t.Variable.String()
		if err != nil {
			return "", err
		}
		return patternDef.Transformer(value), nil
	}

	// Process default value logic first, regardless of KeepUnset setting
	if t.ExpType >= itemPlus && t.Default != nil {
		switch t.ExpType {
		case itemColonDash, itemColonEquals:
			// For colon operators, check if variable is set AND not empty
			if t.Variable.isSet() && t.Variable.Env.Get(t.Variable.Ident) != "" {
				return t.Variable.String()
			}
			return t.Default.String()
		case itemPlus:
			// + operator: return alternate if variable is set (regardless of value)
			if t.Variable.isSet() {
				return t.Default.String()
			}
			return "", nil
		case itemColonPlus:
			// :+ operator: return alternate if variable is set AND not empty
			if t.Variable.isSet() && t.Variable.Env.Get(t.Variable.Ident) != "" {
				return t.Default.String()
			}
			return "", nil
		default:
			// For non-colon operators (dash, equals), check if variable is set
			if !t.Variable.isSet() {
				return t.Default.String()
			}
		}
	}

	// If KeepUnset is enabled and variable is not set, return source text
	// (only if no defaults were processed above)
	if t.Variable.Restrict.KeepUnset && !t.Variable.isSet() {
		// Construct the source text format from ident
		return "${" + t.Variable.Ident + "}", nil
	}

	return t.Variable.String()
}
