package client

import (
	"io"
	"net/http"
)

func RunPythonHelloWorld() {
	resp, err := http.Get("http://controller-python:5000/")
	if err != nil {
		// println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				// println(err->)
			}
			bodyString := string(bodyBytes)
			println(bodyString)
		}
	}
}
