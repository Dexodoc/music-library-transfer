package main

import (
	"fmt"
	"os"
	"strconv"
)


func main() {

	os.Setenv("SPOTIFY_ID", "c8a2b5d8d88449f5b936603275a4f3fe")
	os.Setenv("SPOTIFY_SECRET ", "afc789be90034389b8a0ab2d0a845aa2")

	if false{appleDeveloperToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlRSVUM3NzdEWkwifQ.eyJpYXQiOjE2NTgzMzYzNzksImV4cCI6MTY3Mzg4ODM3OSwiaXNzIjoiM1M1OVVGV1A4WSJ9.ywmXqlzN700OszGHLXePHOk2ifJ3EFWk_mF5XZT8w9bPgK9epenRO62tbpjdVMNgu6p6Hio6095Wig73_3xtwA"
	
	// Currently personal token
	musicUserToken := "AmIKSqcLaVo8ziS52+UVu6Qf+OW8vg8qXAhbpRi/+eh0ObM4AZ22pTVIZdxhsbpHCyiIOtjbmA8jKRkEqQL3+yo9AQSAJgTzu9l1CjSY1WLPOvl71B5bi9ZpFE/FhKahuqyP9Lf+aHV9H3bhThaPrWTrLO1HjxMbaxWjTzgoJVzGqdbviQdfYzxcHTe4nAzBNxKHysAm00SGl8xaVviEzb4Gs5sN1cpRrbso2TjO/cvlfvX+9g=="

	totalSongCount := getSongCount(appleDeveloperToken, musicUserToken)
    f, _ := os.Create("songList.txt") 
	
    defer f.Close()

    for offset := 0; offset < totalSongCount; offset += 100 {
		fmt.Println(offset)

		resData := makeSongRequest(appleDeveloperToken, musicUserToken, 100, offset)
		
        for i, v := range resData.Data {
            f.WriteString(strconv.Itoa(offset + i) + " " + v.Attributes.ArtistName + " - " + v.Attributes.Name + "\n")
        }

	}}

	runSpotify()
}
