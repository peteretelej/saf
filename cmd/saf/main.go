package main

import (
	"flag"
	"fmt"
	"github.com/peteretelej/saf"
	"os"
)

var (
	bundles = flag.NewFlagSet("bundles", flag.ExitOnError)
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "bundles")
	}
	switch os.Args[1] {
	case "bundles":
		b, err := saf.GetBundles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get bundles: %v\n", err)
			os.Exit(1)
		}
		b.PrintTo(os.Stdout)
	}

}
