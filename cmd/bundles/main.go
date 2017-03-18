package main

import (
	"github.com/fatih/color"
	"github.com/peteretelej/saf"
)

func main() {
	b, err := saf.GetBundles()
	if err != nil {
		color.Red("Failed to get bundles: %v", err)
		return
	}

	c := color.New(color.FgGreen).Add(color.Underline)
	c.Println("SAFARICOM BALANCE")
	//	color.White("Safcom Line: 0%s", b.Line)
	c = color.New(color.FgGreen).Add(color.Bold)
	c.Printf("Data Bundles: %s\n", b.Bundles)
	c = color.New(color.FgYellow).Add(color.Bold)
	c.Printf("Airtime Bal: %.2f /=\n", b.Airtime)
}
