package spdx

import (
	"fmt"
	"io"
	"strings"
)

type Token string

const (
	UNSET Token = ""
	AND         = "AND"
	OR          = "OR"
	WITH        = "WITH"
)

type LicenseID string

func NewLicenseID(s string) (LicenseID, error) {
	needle := s
	if strings.HasSuffix(s, "+") {
		needle = s[:len(s)-1]
	}
	for _, known := range ALL {
		if needle == known {
			return LicenseID(s), nil
		}
	}
	return "", fmt.Errorf("unknown license: %s", s)
}

type LicenseExceptionID string

func NewLicenseExceptionID(s string) (LicenseExceptionID, error) {
	for _, known := range LicenseExceptions {
		if s == known {
			return LicenseExceptionID(s), nil
		}
	}
	return "", fmt.Errorf("unknown license exception: %s", s)
}

type CompoundExpr struct {
	Compound   []*CompoundExpr
	CompoundOP Token

	Simple    []LicenseID
	Exception LicenseExceptionID
}

type Parser struct {
	s *Scanner

	last string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (*CompoundExpr, error) {
	return p.parse(0)
}

func (p *Parser) parse(depth int) (*CompoundExpr, error) {
	// TODO: implement precedence of a <license-expression>:
	//    +
	//    WITH
	//    AND
	//    OR
	cur := &CompoundExpr{}
	for p.s.Scan() {
		tok := p.s.Text()
		if p.last == WITH {
			ex, err := NewLicenseExceptionID(tok)
			if err != nil {
				return nil, err
			}
			cur.Exception = ex
			p.last = string(ex)
			continue
		}

		switch tok {
		case "(":
			new, err := p.parse(depth + 1)
			if err != nil {
				return cur, err
			}
			cur.Compound = append(cur.Compound, new)
		case ")":
			if depth == 0 {
				return nil, fmt.Errorf("unbalanced parenthesis")
			}
			if len(cur.Simple) > 1 && cur.CompoundOP == UNSET {
				return nil, fmt.Errorf("need operator in %v", cur.Simple)
			}

			return cur, nil
		case AND:
			if cur.CompoundOP != UNSET && cur.CompoundOP != AND {
				return nil, fmt.Errorf("inconsistent operator in %v, was %v before, changes to %v", cur, cur.CompoundOP, "AND")
			}
			cur.CompoundOP = AND
		case OR:
			if cur.CompoundOP != UNSET && cur.CompoundOP != OR {
				return nil, fmt.Errorf("inconsistent operator in %v, was %v before, changes to %v", cur, cur.CompoundOP, "OR")
			}
			cur.CompoundOP = OR
		case WITH:
			cur.Exception = WITH
		default:
			id, err := NewLicenseID(tok)
			if err != nil {
				return nil, err
			}
			cur.Simple = append(cur.Simple, id)
		}
		p.last = tok
	}
	if err := p.s.Err(); err != nil {
		return nil, err
	}

	if len(cur.Simple) > 1 && cur.CompoundOP == UNSET {
		return nil, fmt.Errorf("need operator in %v", cur.Simple)
	}

	return cur, nil
}
