package parse

import (
	"fmt"
)

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
