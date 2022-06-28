package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get[T any](uri string, queries map[string]string, headers map[string]string) (T, error) {
	var response T
	query := ""

	if len(queries) > 0 {
		for key, value := range queries {
			var startChar string
			if query == "" {
				startChar = "?"
			} else {
				startChar = "&"
			}

			query = fmt.Sprintf("%s%s=%s", startChar, key, value)
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", uri, url.QueryEscape(query)), nil)
	if err != nil {
		return response, errors.New("Could not create request")
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return response, errors.New("Network error")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, errors.New("Reading body failed")
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return response, errors.New("Couldn't map response to generic type T")
	}

	return response, nil
}
