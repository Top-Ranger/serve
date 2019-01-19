// SPDX-License-Identifier: Apache-2.0
// Copyright 2017,2018,2019 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
