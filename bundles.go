package saf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Public http endpoints for getting required data
const (
	lineURL    = "http://www.safaricom.com/bundles/js/get.jsp"
	bundlesURL = "http://www.safaricom.com/bundles/GetSubDetails"
)

// Bundles defines the structure for bundles values returned
type Bundles struct {
	Line    string
	AccType string
	Bundles string
	Airtime float64
}

// GetBundles returns the bundles for the line in use
func GetBundles() (*Bundles, error) {
	l, err := line()
	if err != nil {
		return nil, fmt.Errorf("Unable to get line: %v", err)
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
ParseBundles parses raw bundles html to Bundles value, especially useful in processing data sent from clients (e.g. xhr)

Example raw bundles html:
	<table border=0 id='dataAccount'>
	<tr><td><span>Account Types</span></td><td>Prepaid</td></tr><tr><td><span>Data Bundle</span></td><td>8056.47 MBs</td></tr><tr><td><span>Data Bundle Expiry Date </span></td><td>31-Dec-2036 19:00</td></tr><tr><td><span>Airtime Balances</span></td><td>0.03</td></tr><tr><td><span>Airtime Expiry Date </span></td><td>01-Jan-2037 00:00</td></tr>
	<tr><td colspan=2><div id='refresh'><a href='#'>Refresh</div></td></tr>
	<table>
*/
func ParseBundles(rawBundles string) (*Bundles, error) {
	rows := strings.Split(rawBundles, "<tr>")

	res := map[string]string{"Account Types": "", "Data Bundle</span>": "", "Airtime Balances": ""}

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
	for _, val := range res {
		if val == "" {
			return nil, errors.New("failed to parse bundles")
		}
	}
	airtime, err := strconv.ParseFloat(res["Airtime Balances"], 64)
	if err != nil {
		return nil, errors.New("failed to parse bundles: unable to get airtime")
	}
	return &Bundles{
		AccType: res["Account Types"],
		Bundles: res["Data Bundle</span>"],
		Airtime: airtime,
	}, nil

}

// line simply returns the line, and an error if any
func line() (string, error) {
	resp, err := http.Get(lineURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(dat)), nil
}
