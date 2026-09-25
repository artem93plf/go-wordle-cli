// Package main реализует точку входа приложения «5 букв».
package main

import (
	"fiveletters/cli"
	"flag"
	"fmt"
	"os"
)

func main() {
	seed := flag.Int64("seed", 42, "seed для выбора загаданного слова")
	flag.Parse()

	if err := cli.Run(os.Stdin, os.Stdout, *seed); err != nil {
		if _, printErr := fmt.Fprintln(os.Stderr, err); printErr != nil {
			os.Exit(1)
		}

		os.Exit(1)
	}
}
