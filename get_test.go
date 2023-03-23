package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type Response struct {
	Url string `json:"url"`
}

type URLS struct {
	longUrl  string
	shortUrl string
}

var urls []URLS = []URLS{
	{"github.com", "localhost:8080/0000000001"},
	{"instagram.com", "localhost:8080/0000000002"},
	{"spotify.com", "localhost:8080/0000000003"},
	{"goodreads.com", "localhost:8080/0000000004"},
	{"moodle.com", "localhost:8080/000000000c"},
}

func getJson(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("No response from request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("Bad Request 1")
	}
	var responseObject Response
	json.Unmarshal(body, &responseObject)
	return responseObject.Url, nil
}

func TestGet(t *testing.T) {

	for i := range urls {
		fullUrl := "http://" + urls[i].shortUrl

		str, err := getJson(fullUrl)

		if err != nil {
			t.Errorf("Error parsing json. Error: %s. Url: %s", err.Error(), fullUrl)
			continue
		}
		if str != urls[i].longUrl {
			t.Errorf("Expecred: '%s', but the output is: '%s'", urls[i].longUrl, str)
		}
	}
}
