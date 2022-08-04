package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Structs for destructuring response data
type LibraryResponse struct {
	Next string
	Data []struct {
		Attributes struct {
			ArtistName string
			Name       string
		}
	}
	Meta struct {
		Total int
	}
}

type StorefrontResponse struct {
	Data []struct {
		ID string
	}
}

type SearchResponse struct {
	Results struct {
		Songs struct {
			Data []struct {
				ID         string
				Attributes struct {
					AlbumName  string
					Name       string
					ArtistName string
				}
			}
		}
	}
}

func getMusicUserToken() string {
	response := make(chan string)
	go musicToken(response)
	resp := <-response
	return resp
}

func musicToken(res chan string) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":8000", Handler: m}
	fmt.Println("Please log in to Apple Music by visiting the following page in your browser: http://localhost:8000")
	m.HandleFunc("/return", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		token = strings.ReplaceAll(token, " ", "+")
		res <- token
		s.Shutdown(context.Background())
	})
	fs := http.FileServer(http.Dir("./static"))
	m.HandleFunc("/", fs.ServeHTTP)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getSongCountApple(devToken string, userToken string) int {
	return getSongsApple(devToken, userToken, 1, 0).Meta.Total
}

func getSongsApple(devToken string, userToken string, limit int, offset int) (r LibraryResponse) {
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
	fmt.Println("Getting all songs from Apple Music")

	songs := []Song{}
	songList := SongList{Songs: songs}

	songCount := getSongCountApple(devToken, userToken)
	fmt.Printf("Getting a total of %d songs \n", songCount)

	for offset := 0; offset < songCount; offset += 100 {
		resData := getSongsApple(devToken, userToken, 100, offset)
		for _, v := range resData.Data {
			newSong := Song{name: v.Attributes.Name, artists: []string{v.Attributes.ArtistName}}
			songList.AddItem(newSong)
		}
		upper := offset+100
		if upper >= songCount {
			upper = songCount
		}
		fmt.Printf("Got songs %d - %d \n", offset, upper)
	}
	return songList
}

func getStorefrontApple(devToken string, userToken string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.music.apple.com/v1/me/storefront", nil)

	req.Header.Set("Authorization", "Bearer "+devToken)
	req.Header.Set("Music-User-Token", userToken)
	r := StorefrontResponse{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &r)
	if(len(r.Data) == 0){
		return ""
	}
	return r.Data[0].ID
}

func getSongIdApple(devToken string, userToken string, song Song) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.music.apple.com/v1/catalog/"+getStorefrontApple(devToken, userToken)+"/search", nil)

	// Setting headers for user and developer authentication
	req.Header.Set("Authorization", "Bearer "+devToken)
	req.Header.Set("Music-User-Token", userToken)
	r := SearchResponse{}
	q := req.URL.Query()
	q.Add("types", "songs")
	q.Add("term", song.name+" "+song.artists[0])
	req.URL.RawQuery = q.Encode()
	fmt.Println("Searching for " + song.name)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &r)
	if(len(r.Results.Songs.Data) == 0){
		return ""
	}
	return r.Results.Songs.Data[0].ID
}

func addAllSongsApple(devToken string, userToken string, songs SongList) {
	fmt.Printf("Attempting to add %d songs to Apple Music Library \n", len(songs.Songs))

	ids := make([]string, len(songs.Songs))

	fmt.Println("Searching for songs in Apple Music Library")
	for i, v := range songs.Songs {
		ids[i] = getSongIdApple(devToken, userToken, v)
	}
	fmt.Println("Adding songs to Apple Music Library")
	for i := 0; i < len(songs.Songs); i += 10 {
		upper := i+10
		if upper >= len(songs.Songs) {
			upper = len(songs.Songs)
		} 
		addSongIdApple(devToken, userToken, ids[i:upper])
		fmt.Printf("Added song %d - %d to library \n", i, upper)
	}
}

func addSongIdApple(devToken string, userToken string, ids []string) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.music.apple.com/v1/me/library", nil)

	req.Header.Set("Authorization", "Bearer "+devToken)
	req.Header.Set("Music-User-Token", userToken)

	q := req.URL.Query()
	q.Add("ids[songs]", strings.Join(ids, ","))
	req.URL.RawQuery = q.Encode()

	client.Do(req)
}
