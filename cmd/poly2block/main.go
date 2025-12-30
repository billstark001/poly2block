package main

import (
	"fmt"
	"os"

	"github.com/billstark001/poly2block/cmd/poly2block/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
