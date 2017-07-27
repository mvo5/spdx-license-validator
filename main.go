package main

import (
	"bufio"
	"fmt"
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

type compoundExpr struct {
	simple     licenseID
	compound   []licenseID
	cmmpoundOP int
}

func spdxSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//fmt.Printf("%v %q %q %v\n", len(data), string(data), string(data[0]), atEOF)
	// skip WS
	start := 0
	for ; start < len(data); start++ {
		if data[start] != ' ' {
			break
		}
	}
	if start == len(data) {
		return start, nil, nil
	}

	// found (,)
	switch data[start] {
	case '(', ')':
		return start + 1, data[start : start+1], nil
	}

	// found non-ws, non-(), must be a token
	for i := start; i < len(data); i++ {
		switch data[i] {
		// token finished
		case ' ', '\n':
			return i + 1, data[start:i], nil
			// found (,) - we need to rescan it
		case '(', ')':
			return i, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(spdxSplit)
	for scanner.Scan() {
		fmt.Printf("%q\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Invalid input: %s", err)
	}
}
