package main

import (
	"fmt"
	"os"
)

type SongList struct {
	Songs []Song
}

type Song struct {
	name    string
	artists []string
}

func (sl *SongList) AddItem(song Song) []Song {
	sl.Songs = append(sl.Songs, song)
	return sl.Songs
}

func main() {
	os.Setenv("SPOTIFY_ID", "c8a2b5d8d88449f5b936603275a4f3fe")
	os.Setenv("SPOTIFY_SECRET", "afc789be90034389b8a0ab2d0a845aa2")
	os.Setenv("APPLE_DEV_TOKEN", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlRSVUM3NzdEWkwifQ.eyJpYXQiOjE2NTkzOTYxMDMsImV4cCI6MTY3NDk0ODEwMywiaXNzIjoiM1M1OVVGV1A4WSJ9.Z5UJwcOyil_2DUjqmCP1ztCYCNnfa3RpwkYrf5rmbJEcgf0iNuzYOgltHUZdtASsxBlYVoI5be3ZW2OchAEEBg")

	
	client, _ := NewSpotifyClientBuilder().GetClient()
	songs := getAllSongsSpotify(*client)

	pain := getMusicUserToken()
	addAllSongsApple(os.Getenv("APPLE_DEV_TOKEN"), pain, songs)
	// song := Song{name: "Luck Be a Lady", artists: []string{"Frank Sinatra"}}
	// song2 := Song{name: "We Don't Know", artists: []string{"The Strumbellas"}}
	// song3 := Song{name: "Spirits", artists: []string{"The Strumbellas"}}
	// song4 := Song{name: "Just Hold On", artists: []string{"Steve Aoki"}}
	// song5 := Song{name: "Too Young to Burn", artists: []string{"Sonny & The Sunsets"}}
	// songs := SongList{Songs: []Song{song, song2, song3, song4, song5}}
	// addAllSongsSpotify(*client, songs)
	// // addSongIdSpotify(*client, id)
	fmt.Println("ADAW")
}
