package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func postQuery(url string, query string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"url": query,
	})

	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "x-www-form-urlencoded", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var responseObject postUrl
	json.Unmarshal(body, &responseObject)

	return responseObject.Url, nil
}

func TestPost(t *testing.T) {

	for i := range urls {

		res, err := postQuery("http://localhost:8080/post", urls[i].longUrl)

		if err != nil {
			t.Errorf("1 %s", err.Error())
		}

		if res != urls[i].shortUrl {
			t.Errorf("The %s was expected. The %s got caught", urls[i].shortUrl, res)
		}

	}

}
