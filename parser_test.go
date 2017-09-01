package spdx_test

import (
	"bytes"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/mvo5/spdx-license-validator"
)

func Test(t *testing.T) { TestingT(t) }

type spdxSuite struct{}

var _ = Suite(&spdxSuite{})

func (s *spdxSuite) TestParseHappy(c *C) {
	for _, t := range []string{
		"GPL-2.0",
		"GPL-2.0+",
		"GPL-2.0 AND BSD-2-Clause",
		"GPL-2.0 OR BSD-2-Clause",
		"GPL-2.0 WITH GCC-exception-3.1",
		"(GPL-2.0 AND BSD-2-Clause)",
		"GPL-2.0 AND (BSD-2-Clause OR 0BSD)",
		"GPL-2.0 AND (BSD-2-Clause OR 0BSD) WITH GCC-exception-3.1",
		"((GPL-2.0 AND (BSD-2-Clause OR 0BSD)) OR GPL-3.0) ",
	} {
		parser := spdx.NewParser(bytes.NewBufferString(t))
		err := parser.Validate()
		c.Check(err, IsNil, Commentf("input: %q", t))
	}
}

func (s *spdxSuite) TestParseError(c *C) {
	for _, t := range []struct {
		inp    string
		errStr string
	}{
		{"FOO", `unknown license: FOO`},
		{"GPL-2.0 GPL-3.0", `unexpected token: "GPL-3.0"`},
		{"(GPL-2.0))", "unbalanced parenthesis"},
		{"(GPL-2.0", `expected closing parenthesis got ""`},
		{"OR", "expected left license with operator OR"},
		{"OR GPL-2.0", "expected left license with operator OR"},
		{"GPL-2.0 OR", "expected right license with operator OR"},
		{"GPL-2.0 WITH BAR", "unknown license exception: BAR"},
		{"GPL-2.0 WITH (foo)", `unknown license exception: foo`},
	} {
		parser := spdx.NewParser(bytes.NewBufferString(t.inp))
		err := parser.Validate()
		c.Check(err, ErrorMatches, t.errStr, Commentf("input: %q", t.inp))
	}
}
