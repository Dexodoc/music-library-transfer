package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Response struct {
	Next string
	Data []Song
	Meta Meta
}

type Meta struct {
	Total int
}

type Song struct {
	Attributes Attributes
}

type Attributes struct {
	ArtistName string
	Name       string
}

func getSongCount(devToken string, userToken string) int {
	return makeSongRequest(devToken, userToken, 100, 0).Meta.Total
}

func makeSongRequest(devToken string, userToken string, limit int, offset int) (r Response) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.music.apple.com/v1/me/library/songs", nil)

	// Setting headers for user and developer authentication
	req.Header.Set("Authorization", "Bearer "+devToken)
	req.Header.Set("Music-User-Token", userToken)

	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	req.URL.RawQuery = q.Encode()

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &r)
	
	return
}

