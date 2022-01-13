package utils

import (
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
