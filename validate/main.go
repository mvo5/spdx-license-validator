package main

import (
	"log"
	"os"

	"github.com/mvo5/spdx-license-validator"
)

func output(cur *spdx.CompoundExpr, depth int) {
	for _, comp := range cur.Compound {
		output(comp, depth+1)
	}
	for i := 0; i < depth; i++ {
		print(" ")
	}
	print(cur.CompoundOP)
	print(" ")

	for _, s := range cur.Simple {
		print(s)
		print(" ")
	}
	println()
}

func main() {
	parser := spdx.NewParser(os.Stdin)
	res, err := parser.Parse()
	if err != nil {
		log.Fatalf("cannot parse: %s", err)
	}
	output(res, 0)
}
