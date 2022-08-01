package main

import (
	"fmt"
	"os"
)

type SongList struct {
	Songs []Song
}

type Song struct {
	name string
	artists []string
}

func (sl *SongList) AddItem(song Song) []Song{
	sl.Songs = append(sl.Songs, song)
	return sl.Songs
}


func main() {
	os.Setenv("SPOTIFY_ID", "c8a2b5d8d88449f5b936603275a4f3fe")
	os.Setenv("SPOTIFY_SECRET", "afc789be90034389b8a0ab2d0a845aa2")
	os.Setenv("APPLE_DEV_TOKEN", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlRSVUM3NzdEWkwifQ.eyJpYXQiOjE2NTkzOTYxMDMsImV4cCI6MTY3NDk0ODEwMywiaXNzIjoiM1M1OVVGV1A4WSJ9.Z5UJwcOyil_2DUjqmCP1ztCYCNnfa3RpwkYrf5rmbJEcgf0iNuzYOgltHUZdtASsxBlYVoI5be3ZW2OchAEEBg")
	
	pain := getMusicUserToken()
	fmt.Println(pain)
	fmt.Println(getSongCountApple(os.Getenv("APPLE_DEV_TOKEN"), pain))
}
