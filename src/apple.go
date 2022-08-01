package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Structs for destructuring response data

type Response struct {
	Next string
	Data []AppleSong
	Meta AppleMeta
}

type AppleMeta struct {
	Total int
}

type AppleSong struct {
	Attributes AppleAttributes
}

type AppleAttributes struct {
	ArtistName string
	Name       string
}

func getSongCountApple(devToken string, userToken string) int {
	return getSongsApple(devToken, userToken, 100, 0).Meta.Total
}

func getSongsApple(devToken string, userToken string, limit int, offset int) (r Response) {
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

func getAllSongsApple(devToken string, userToken string) SongList {
	songs := []Song{}
	songList := SongList{Songs: songs}
	
	songCount := getSongCountApple(devToken, userToken)

	for offset := 0; offset < songCount; offset += 100 {
		resData := getSongsApple(devToken, userToken, 100, offset)
		
		for _, v := range resData.Data {
			newSong := Song{name: v.Attributes.Name, artists: []string{v.Attributes.ArtistName}}
			songList.AddItem(newSong)
		}
	}

	return songList
}

