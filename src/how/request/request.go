package request

import (
	"net/http"
)

func Get(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	return client.Do(req)
}
