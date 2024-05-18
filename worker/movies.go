package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func getMovies(baseApiUrl, apiKey string, genres []int) ([]byte, error) {
	url := fmt.Sprintf("%s&api_key=%s&with_genres=", baseApiUrl, apiKey)
	for _, genre := range genres {
		url = url + fmt.Sprintf("%d|", genre)
	}
	var err error
	for i := 0; i < 3; i++ {
		url = strings.TrimSuffix(url, "|")
		req, err := http.NewRequest("GET", url, nil)
		if nil == err {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				continue
			}
			return body, nil
		}
		time.Sleep(5 * time.Second)
	}
	return nil, err
}
