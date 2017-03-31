package saf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Public http endpoints for getting required data
const (
	lineURL    = "http://www.safaricom.com/bundles/js/get.jsp"
	bundlesURL = "http://www.safaricom.com/bundles/GetSubDetails"
)

// Bundles defines the structure for bundles values returned
type Bundles struct {
	Line          string
	AccType       string
	Bundles       string
	BundlesExpiry string
	Airtime       float64
	BongaSMS      string
	BongaBalance  string
}

// PrettyPrint prints bundles information to stdout
func (b *Bundles) PrettyPrint() {
	var (
		head   = color.New(color.FgGreen, color.Bold).Add(color.Underline).SprintFunc()
		info   = color.New(color.FgWhite, color.BgGreen).SprintFunc()
		red    = color.New(color.FgRed).SprintFunc()
		green  = color.New(color.FgGreen).SprintFunc()
		yellow = color.New(color.FgYellow).SprintFunc()
	)

	fmt.Printf("\n%s\n", head("SAFARICOM BALANCE"))
	if b.Line != "" && len(b.Line) < 30 {
		fmt.Printf("Safcom Line: %s\n", info("0"+b.Line))
	}
	airtime := fmt.Sprintf("%.2f", b.Airtime)
	fmt.Printf("Airtime Balance: %s/=\n", (airtime))
	fmt.Printf("Data Bundles: %s\n", green(b.Bundles))
	fmt.Printf("Bundles Expiry: %s\n\n", red(b.BundlesExpiry))
	if b.BongaSMS != "" {
		fmt.Printf("Bonga SMS: %s\n", yellow(b.BongaSMS))
	}
	if b.BongaBalance != "" {
		fmt.Printf("Bonga Balance: %s\n", yellow(b.BongaBalance))
	}
}

// GetBundles returns the bundles for the line in use
func GetBundles() (*Bundles, error) {
	l, err := line()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(bundlesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	b, err := ParseBundles(string(dat))
	if err != nil {
		return nil, err
	}
	b.Line = l
	return b, nil
}

/*
ParseBundles parses raw bundles html to Bundles value.
Raw bundles html are as received from the Saf GetSubDetails http endpoint.
Especially useful in processing data sent from clients (e.g. xhr)
*/
func ParseBundles(rawBundles string) (*Bundles, error) {
	rows := strings.Split(rawBundles, "<tr>")

	res := map[string]string{
		"Account Types":           "",
		"Data Bundle</span>":      "",
		"Data Bundle Expiry Date": "",
		"Airtime Balances":        "",
		"Bonga SMS</":             "",
		"Bonga Balance":           "",
	}

	for _, val := range rows {
		if !strings.Contains(val, "</tr>") {
			continue
		}
		var title string
		for key := range res {
			if strings.Contains(val, key) {
				title = key
			}
		}
		if title == "" {
			continue
		}

		if strings.Count(val, "<td>") < 2 {
			continue
		}
		val2 := strings.Split(val, "<td>")[2]
		val2 = strings.Replace(val2, "<td>", "", 1)
		val2 = strings.Replace(val2, "</td></tr>", "", 1)
		res[title] = val2
	}
	var ok bool
	for _, val := range res {
		if val != "" {
			ok = true // some data found
			break
		}
	}
	if !ok {
		return nil, errors.New("Unable to get bundles")
	}

	airtime, err := strconv.ParseFloat(res["Airtime Balances"], 64)
	if err != nil {
		return nil, errors.New("failed to parse bundles: unable to get airtime")
	}
	return &Bundles{
		AccType:       res["Account Types"],
		Bundles:       res["Data Bundle</span>"],
		Airtime:       airtime,
		BundlesExpiry: res["Data Bundle Expiry Date"],
		BongaSMS:      res["Bonga SMS</"],
		BongaBalance:  res["Bonga Balance"],
	}, nil

}

// line simply returns the line, and an error if any
func line() (string, error) {
	resp, err := http.Get(lineURL)
	if err != nil {
		return "", errors.New("no internet connection")
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	l := strings.TrimSpace(string(dat))
	if l == "" {
		return "", errors.New("you aren't connected through Safaricom")
	}
	return l, nil
}
