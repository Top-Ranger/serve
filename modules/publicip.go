package modules

import (
	"io/ioutil"
	"net/http"
)

const ipGetAddress = "https://api.ipify.org/"

func GetPublicIp() (string, error) {
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
