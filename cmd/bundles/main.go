package main

import (
	"fmt"
	"os"

	"github.com/peteretelej/saf"
)

func main() {
	b, err := saf.GetBundles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get bundles: %v\n", err)
		os.Exit(1)
	}
	b.PrettyPrint()

}
