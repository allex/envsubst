// Most of the code in this package taken from golang/text/template/parse
package parse

import (
	"errors"
	"strings"
)

// A mode value is a set of flags (or 0). They control parser behavior.
type Mode int

// Mode for parser behaviour
const (
	Quick     Mode = iota // stop parsing after first error encoutered and return
	AllErrors             // report all errors
)

// Restrictions controls the parsing and substitution behavior of environment variables.
// These options determine how the parser handles undefined variables, empty variables,
// numeric variables, and variable matching patterns.
type Restrictions struct {
	// NoUnset when true causes the parser to return an error if a variable is not set.
	// When false (default), unset variables are substituted with empty strings.
	// Example: ${UNDEFINED_VAR} will cause an error if NoUnset is true.
	NoUnset bool

	// NoEmpty when true causes the parser to return an error if a variable is set but empty.
	// When false (default), empty variables are substituted with empty strings.
	// Example: If VAR="" then ${VAR} will cause an error if NoEmpty is true.
	NoEmpty bool

	// NoDigit when true causes the parser to ignore variables that start with a digit.
	// When false (default), numeric variables are processed normally.
	// Example: $1 and ${1} will be treated as literal text if NoDigit is true.
	NoDigit bool

	// KeepUnset when true causes undefined variables to be kept as their original text
	// instead of being substituted with empty strings or causing errors.
	// When true, this option automatically disables NoUnset and NoEmpty restrictions.
	// Example: ${UNDEFINED_VAR} will remain as "${UNDEFINED_VAR}" in the output.
	KeepUnset bool

	// VarMatcher is an optional predicate function to filter valid variable tokens.
	// If provided, only variables that pass this filter will be processed.
	// Variables that don't match will be treated as literal text.
	VarMatcher varMatcher
}

// Parser type initializer
type Parser struct {
	Name     string // name of the processing template
	Env      *Env
	Restrict *Restrictions
	Mode     Mode
	// parsing state;
	lex       *lexer
	token     [3]item // three-token lookahead
	peekCount int
	nodes     []Node
}

// New allocates a new Parser with the given name.
func New(name string, env *Env, r *Restrictions) *Parser {
	if r != nil && r.KeepUnset {
		r.NoEmpty = false
		r.NoUnset = false
	}
	return &Parser{
		Name:     name,
		Env:      env,
		Restrict: r,
	}
}

// Parse parses the given string.
func (p *Parser) Parse(text string) (string, error) {
	p.lex = lex(text, p.Restrict.NoDigit, p.Restrict.VarMatcher)
	// Build internal array of all unset or empty vars here
	var errs []error
	// clean parse state
	p.nodes = make([]Node, 0)
	p.peekCount = 0
	if err := p.parse(); err != nil {
		switch p.Mode {
		case Quick:
			return "", err
		case AllErrors:
			errs = append(errs, err)
		}
	}
	var out string
	for _, node := range p.nodes {
		s, err := node.String()
		if err != nil {
			switch p.Mode {
			case Quick:
				return "", err
			case AllErrors:
				errs = append(errs, err)
			}
		}
		out += s
	}
	if len(errs) > 0 {
		var b strings.Builder
		for i, err := range errs {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(err.Error())
		}
		return "", errors.New(b.String())
	}
	return out, nil
}

// parse is the top-level parser for the template.
// It runs to EOF and return an error if something isn't right.
func (p *Parser) parse() error {
Loop:
	for {
		switch t := p.next(); t.typ {
		case itemEOF:
			break Loop
		case itemError:
			return p.errorf(t.val)
		case itemVariable:
			varNode := NewVariableWithSource(strings.TrimPrefix(t.val, "$"), t.val, p.Env, p.Restrict)
			p.nodes = append(p.nodes, varNode)
		case itemLeftDelim:
			if p.peek().typ == itemVariable {
				n, err := p.action()
				if err != nil {
					return err
				}
				p.nodes = append(p.nodes, n)
				continue
			}
			fallthrough
		default:
			textNode := NewText(t.val)
			p.nodes = append(p.nodes, textNode)
		}
	}
	return nil
}

// Parse substitution. first item is a variable.
func (p *Parser) action() (Node, error) {
	var expType itemType
	var defaultNode Node

	varToken := p.next()
	varNode := NewVariable(varToken.val, p.Env, p.Restrict)

	// Build source text for the substitution - start with the basic form
	sourceText := "${" + varToken.val

Loop:
	for {
		switch t := p.next(); t.typ {
		case itemRightDelim:
			sourceText += "}"
			break Loop
		case itemError:
			return nil, p.errorf(t.val)
		case itemVariable:
			defaultNode = NewVariable(strings.TrimPrefix(t.val, "$"), p.Env, p.Restrict)
			sourceText += t.val
		case itemText:
			n := NewText(t.val)
			sourceText += t.val
		Text:
			for {
				switch p.peek().typ {
				case itemRightDelim, itemError, itemEOF:
					break Text
				default:
					// patch to accept all kind of chars
					nextToken := p.next()
					n.Text += nextToken.val
					sourceText += nextToken.val
				}
			}
			defaultNode = n
		default:
			expType = t.typ
			// Add operator tokens to source text
			switch t.typ {
			case itemColonDash:
				sourceText += ":-"
			case itemColonEquals:
				sourceText += ":="
			case itemColonPlus:
				sourceText += ":+"
			case itemDash:
				sourceText += "-"
			case itemEquals:
				sourceText += "="
			case itemPlus:
				sourceText += "+"
			}
		}
	}

	// Update the variable node with source text
	varNode.Source = sourceText

	return &SubstitutionNode{NodeSubstitution, expType, varNode, defaultNode, sourceText}, nil
}

func (p *Parser) errorf(s string) error {
	return errors.New(s)
}

// next returns the next token.
func (p *Parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	return p.token[p.peekCount]
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.peekCount++
}

// peek returns but does not consume the next token.
func (p *Parser) peek() item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.nextItem()
	return p.token[0]
}
