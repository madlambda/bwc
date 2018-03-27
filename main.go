package main

import (
	"fmt"
	"io"
	"os"

	"github.com/madlambda/bwc/infix"
)

func main() {
	for {
		var line [256]byte
		n, err := os.Stdin.Read(line[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s\n", err)
			os.Exit(1)
		}

		tree, err := infix.Parse(string(line[:n-1]))
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		res := infix.Eval(tree)
		fmt.Printf("bin: %b\n", res)
		fmt.Printf("hex: %x\n", res)
	}
}