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
	bundlesURL = "http://www.safaricom.com/bundles/GetSubDetails"
)

// Bundles defines the structure for bundles values returned
type Bundles struct {
	AccType       string
	Bundles       string
	BundlesExpiry string
	Airtime       float64
	BongaSMS      string
	BongaBalance  string
}

func (b Bundles) String() string {
	var (
		head   = color.New(color.FgGreen, color.Bold).Add(color.Underline).SprintFunc()
		red    = color.New(color.FgRed).SprintFunc()
		green  = color.New(color.FgGreen).SprintFunc()
		yellow = color.New(color.FgYellow).SprintFunc()
	)

	bal := head("SAFARICOM BALANCE")
	if b.Airtime != 0 {
		bal += fmt.Sprintf("\nAirtime Balance: %.2f/=", (b.Airtime))
	}
	if b.Bundles != "" {
		bal += fmt.Sprintf("\nData Bundles: %s", green(b.Bundles))
	}
	if b.BundlesExpiry != "" {
		bal += fmt.Sprintf("\nBundles Expiry: %s", red(b.BundlesExpiry))
	}
	if b.BongaSMS != "" {
		bal += fmt.Sprintf("\nBonga SMS: %s", yellow(b.BongaSMS))
	}
	if b.BongaBalance != "" {
		bal += fmt.Sprintf("\nBonga Balance: %s", yellow(b.BongaBalance))
	}
	return bal

}

// GetBundles returns the bundles for the line in use
func GetBundles() (*Bundles, error) {
	resp, err := http.Get(bundlesURL)
	if err != nil {
		return nil, errors.New("no internet connection")
	}
	defer func() { _ = resp.Body.Close() }()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	return ParseBundles(string(dat))
}

/*
ParseBundles parses raw bundles html to Bundles value.
Raw bundles html are as received from the Saf GetSubDetails http endpoint.
Especially useful in processing data sent from clients (e.g. xhr)
*/
func ParseBundles(rawBundles string) (*Bundles, error) {
	if strings.Contains(rawBundles, "not able to capture") {
		return nil, errors.New("you aren't connected through Safaricom")
	}
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
