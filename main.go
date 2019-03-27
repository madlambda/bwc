package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/madlambda/bwc/bwc"
)

var cmd string

func abortonerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		os.Exit(1)
	}
}

func printResult(res int) {
	fmt.Printf("dec: %d\n", res)
	fmt.Printf("bin: %b\n", res)
	fmt.Printf("hex: %x\n", res)
}

func execCmd() (int, error) {
	interp := bwc.NewInterp()

	tree, err := bwc.Parse(cmd)
	if err != nil {
		return 0, err
	}

	res, err := interp.Eval(tree)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func cli() {
	interp := bwc.NewInterp()

	for {
		fmt.Printf("bwc> ")
		var line [256]byte
		n, err := os.Stdin.Read(line[:])
		if err == io.EOF {
			break
		}
		abortonerr(err)

		buf := strings.TrimSpace(string(line[:n-1]))
		if len(buf) == 0 {
			fmt.Printf("\n")
			continue
		}

		tree, err := bwc.Parse(buf)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		res, err := interp.Eval(tree)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}

		printResult(int(res))
	}
}

func main() {
	flag.StringVar(&cmd, "c", "", "Evaluates a command")
	flag.Parse()

	if cmd != "" {
		res, err := execCmd()
		abortonerr(err)
		printResult(res)
		return
	}

	cli()
}
