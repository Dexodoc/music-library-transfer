package main

import (
	"fmt"
	"os"
	"strconv"
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
	// Need to set spotify credentials based on spotify app created in dashboard
	// and ensure redirect uri is et for "http://localhost:8080/callback"
	os.Setenv("SPOTIFY_ID", "")
	os.Setenv("SPOTIFY_SECRET", "")

	// Need to set JWT token generated from apple developer account
	APPLE_DEV_TOKEN := ""

	fmt.Println("Music Library Transfer")

	fmt.Println("Would you like to transfer from Apple Music to Spotify(1) or Spotify to Apple Music(2):")
	direction := getUserInputNumber()
	
	spotifyClient, _ := NewSpotifyClientBuilder().GetClient()
	appleUserToken := getMusicUserToken()

	if direction == 1{
		songs := getAllSongsApple(APPLE_DEV_TOKEN, appleUserToken)
		addAllSongsSpotify(*spotifyClient, songs)
	}else{
		songs := getAllSongsSpotify(*spotifyClient)
		addAllSongsApple(APPLE_DEV_TOKEN, appleUserToken, songs)
	}
}

func getUserInputNumber() int {
	var s string
	var i int
	for {
		_, err := fmt.Scan(&s)
		i, err = strconv.Atoi(s)
		if err != nil || i > 2 || i < 1 {
			fmt.Println("Enter a valid number")
		} else {
			return i
		}
	}
}
