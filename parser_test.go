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
	for _, t := range []struct {
		inp string
		res *spdx.CompoundExpr
	}{
		{
			"GPL-2.0",
			&spdx.CompoundExpr{
				Simple: []spdx.LicenseID{"GPL-2.0"},
			},
		},
		{
			"GPL-2.0+",
			&spdx.CompoundExpr{
				Simple: []spdx.LicenseID{"GPL-2.0+"},
			},
		},
		{
			"GPL-2.0 AND BSD-2-Clause",
			&spdx.CompoundExpr{
				Simple:     []spdx.LicenseID{"GPL-2.0", "BSD-2-Clause"},
				CompoundOP: spdx.AND,
			},
		},
		{
			"GPL-2.0 OR BSD-2-Clause",
			&spdx.CompoundExpr{
				Simple:     []spdx.LicenseID{"GPL-2.0", "BSD-2-Clause"},
				CompoundOP: spdx.OR,
			},
		},
		{
			"GPL-2.0 WITH GCC-exception-3.1",
			&spdx.CompoundExpr{
				Simple:    []spdx.LicenseID{"GPL-2.0"},
				Exception: "GCC-exception-3.1",
			},
		},
		{
			"(GPL-2.0 AND BSD-2-Clause)",
			&spdx.CompoundExpr{
				Compound: []*spdx.CompoundExpr{
					{
						Simple:     []spdx.LicenseID{"GPL-2.0", "BSD-2-Clause"},
						CompoundOP: spdx.AND,
					},
				},
			},
		},
		{
			"GPL-2.0 AND (BSD-2-Clause OR 0BSD)",
			&spdx.CompoundExpr{
				Simple:     []spdx.LicenseID{"GPL-2.0"},
				CompoundOP: spdx.AND,
				Compound: []*spdx.CompoundExpr{
					{
						Simple:     []spdx.LicenseID{"BSD-2-Clause", "0BSD"},
						CompoundOP: spdx.OR,
					},
				},
			},
		},
	} {
		parser := spdx.NewParser(bytes.NewBufferString(t.inp))
		res, err := parser.Parse()
		c.Assert(err, IsNil)
		c.Check(res, DeepEquals, t.res)
	}
}

func (s *spdxSuite) TestParseError(c *C) {
	for _, t := range []struct {
		inp    string
		errStr string
	}{
		{"FOO", "unknown license: FOO"},
		{"(GPL-2.0))", "unbalanced parenthesis"},
		{"GPL-2.0 AND 0BSD OR GPL-3.0", "inconsistent operator .* was AND before, changes to OR"},
		{"GPL-2.0 WITH BAR", "unknown license exception: BAR"},
		{"GPL-2.0 WITH (foo)", `unknown license exception: \(`},
	} {
		parser := spdx.NewParser(bytes.NewBufferString(t.inp))
		_, err := parser.Parse()
		c.Assert(err, ErrorMatches, t.errStr)
	}
}
