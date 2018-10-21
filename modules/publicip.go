// Copyright (c) 2017 Marcus Soll
// SPDX-License-Identifier: MIT

package modules

import (
	"io/ioutil"
	"net/http"
)

// Constant holding the URL for the IP lookup
const ipGetAddress = "https://api.ipify.org/"

// GetPublicIP looks up the current public ip and returns it as a  string.
// The method will block while doing the lookup.
func GetPublicIP() (string, error) {
	resp, err := http.Get(ipGetAddress)
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()

	ip, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
