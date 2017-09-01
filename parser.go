package spdx

import (
	"fmt"
	"io"
	"strings"
)

type Operator string

const (
	UNSET Operator = ""
	AND            = "AND"
	OR             = "OR"
	WITH           = "WITH"
)

func isOperator(tok string) bool {
	return tok == AND || tok == OR || tok == WITH
}

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

type Parser struct {
	s *Scanner

	// state
	last string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Validate() error {
	return p.validate(0)
}

func (p *Parser) advance(id string) error {
	if p.s.Text() != id {
		return fmt.Errorf("expected %q got %q", id, p.s.Text())
	}
	return nil
}

func (p *Parser) validate(depth int) error {
	for p.s.Scan() {
		tok := p.s.Text()

		switch {
		case tok == "(":
			if err := p.validate(depth + 1); err != nil {
				return err
			}
			if p.s.Text() != ")" {
				return fmt.Errorf(`expected closing parenthesis got %q`, p.s.Text())
			}
		case tok == ")":
			if depth == 0 {
				return fmt.Errorf("unbalanced parenthesis")
			}
			return nil
		case isOperator(tok):
			if p.last == "" {
				return fmt.Errorf("expected left license with operator %s", tok)
			}
			if p.last == AND || p.last == OR {
				return fmt.Errorf("expected license-id, got %q", tok)
			}
			if p.last == WITH {
				return fmt.Errorf("expected license-exception-id, got %q", tok)
			}
		default:
			switch {
			case p.last == WITH:
				if _, err := NewLicenseExceptionID(tok); err != nil {
					return err
				}
			case p.last == "", p.last == AND, p.last == OR:
				if _, err := NewLicenseID(tok); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected token: %q", tok)
			}

		}
		p.last = tok
	}
	if err := p.s.Err(); err != nil {
		return err
	}
	if isOperator(p.last) {
		return fmt.Errorf("expected right license with operator %s", p.last)
	}

	return nil
}
