package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/madlambda/bwc/infix"
)

func main() {
	for {
		fmt.Printf("bwc> ")
		var line [256]byte
		n, err := os.Stdin.Read(line[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s\n", err)
			os.Exit(1)
		}

		buf := strings.TrimSpace(string(line[:n-1]))
		if len(buf) == 0 {
			fmt.Printf("\n")
			continue
		}

		tree, err := infix.Parse(buf)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		res, err := infix.Eval(tree)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}

		fmt.Printf("bin: %b\n", int(res))
		fmt.Printf("hex: %x\n", int(res))
	}
}