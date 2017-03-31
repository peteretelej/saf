package main

import (
	"fmt"

	"github.com/peteretelej/saf"
)

func main() {
	b, err := saf.GetBundles()
	if err != nil {
		fmt.Printf("Failed to get bundles: %v\n", err)
		return
	}
	b.PrettyPrint()

}
