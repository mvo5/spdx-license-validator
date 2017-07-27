package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	AND
	OR

	OPENBR    // (
	CLOSINGBR // )
)

type licenseID string

func NewLicenseID(s string) licenseID {
	// FIXME: validate
	return licenseID(s)
}

type compoundExpr struct {
	compound   []*compoundExpr
	compoundOP Token

	simple []licenseID
}

type Parser struct {
	s *Scanner
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (*compoundExpr, error) {
	root := &compoundExpr{}
	cur, err := p.parse(root)
	if cur != root {
		return nil, fmt.Errorf("unbalanced parantheis")
	}
	return cur, err
}

func (p *Parser) parse(cur *compoundExpr) (*compoundExpr, error) {
	for p.s.Scan() {
		tok := p.s.Text()
		switch tok {
		case "(":
			new := &compoundExpr{}
			_, err := p.parse(new)
			if err != nil {
				return cur, err
			}
			cur.compound = append(cur.compound, new)
		case ")":
			return cur, nil
		case "AND", "and":
			if cur.compoundOP != ILLEGAL && cur.compoundOP != AND {
				return nil, fmt.Errorf("cannot chnage op of %v", cur)
			}
			cur.compoundOP = AND
		case "OR", "or":
			if cur.compoundOP != ILLEGAL && cur.compoundOP != OR {
				return nil, fmt.Errorf("cannot chnage op of %v", cur)
			}
			cur.compoundOP = OR
		default:
			cur.simple = append(cur.simple, NewLicenseID(tok))
		}
	}
	if err := p.s.Err(); err != nil {
		return nil, err
	}

	return cur, nil
}

func output(cur *compoundExpr, depth int) {
	for _, comp := range cur.compound {
		output(comp, depth+1)
	}
	for i := 0; i < depth; i++ {
		print(" ")
	}
	print(cur.compoundOP)
	print(" ")

	for _, s := range cur.simple {
		print(s)
		print(" ")
	}
	println()
}

func main() {
	parser := NewParser(os.Stdin)
	res, err := parser.Parse()
	if err != nil {
		log.Fatalf("cannot parse: %s", err)
	}
	output(res, 0)
}
