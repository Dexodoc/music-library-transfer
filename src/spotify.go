package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/zmb3/spotify"
)

type SpotifyClientBuilder struct {
	auth  spotify.Authenticator
	state string
	ch    chan *spotify.Client
}

func NewSpotifyClientBuilder() *SpotifyClientBuilder {
	c := &SpotifyClientBuilder{}

	c.auth = spotify.NewAuthenticator(
		"http://localhost:8080/callback",
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserLibraryRead,
	)

	c.state = randStringBytes(40)
	c.ch = make(chan *spotify.Client)
	return c
}

// GetClient uses the oauth2 flow to get an authenticated Spotify client
func (c *SpotifyClientBuilder) GetClient() (*spotify.Client, error) {
	// start an HTTP server and initiate an oidc flow
	http.HandleFunc("/callback", c.completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":8080", nil)
	url := c.auth.AuthURL(c.state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	client := <-c.ch
	return client, nil
}

func (c *SpotifyClientBuilder) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := c.auth.Token(c.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != c.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, c.state)
	}
	// use the token to get an authenticated client
	client := c.auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	c.ch <- &client
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"

func getAllSongsSpotify(spotifyClient spotify.Client) SongList {
	fmt.Println("Getting all songs from Spotify")

	songs := []Song{}
	songList := SongList{Songs: songs}

	songCount := getSongCountSpotify(spotifyClient)
	fmt.Printf("Getting a total of %d songs \n", songCount)

	for offset := 0; offset < songCount; offset += 50 {
		limit := 50
		options := spotify.Options{Limit: &limit, Offset: &offset}
		songs, err := spotifyClient.CurrentUsersTracksOpt(&options)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range songs.Tracks {
			newSong := Song{name: val.FullTrack.Name, artists: simpleArtistToString(val.Artists)}
			songList.AddItem(newSong)
		}
		upper := offset+50
		if upper >= songCount {
			upper = songCount
		}
		fmt.Printf("Got songs %d - %d \n", offset, upper)
	}

	return songList
}

func getSongCountSpotify(spotifyClient spotify.Client) int {
	tracks, err := spotifyClient.CurrentUsersTracks()
	if err != nil {
		log.Fatal(err)
	}

	return tracks.Total
}

func simpleArtistToString(artists []spotify.SimpleArtist) []string {
	toReturn := make([]string, len(artists))

	for i, e := range artists {
		toReturn[i] = e.Name
	}

	return toReturn
}

func addAllSongsSpotify(spotifyClient spotify.Client, songs SongList) {
	fmt.Printf("Attempting to add %d songs to Spotify Library \n", len(songs.Songs))
	ids := make([]spotify.ID, len(songs.Songs))
	fmt.Println("Searching for songs in Spotify library")
	for i, v := range songs.Songs {
		ids[i] = getSongIdSpotify(spotifyClient, v)
	}

	addSongIdSpotify(spotifyClient, ids)
}

func getSongIdSpotify(spotifyClient spotify.Client, song Song) spotify.ID {
	res, err := spotifyClient.Search(song.name+" "+song.artists[0], spotify.SearchTypeTrack)
	if err != nil {
		fmt.Println(err)
		return spotify.ID("")
	}
	fmt.Println("Searching for " + song.name)
	if len(res.Tracks.Tracks) == 0 {
		return spotify.ID("")
	}
	return res.Tracks.Tracks[0].ID
}

func addSongIdSpotify(spotifyClient spotify.Client, ids []spotify.ID) {
	fmt.Println("Adding songs to Spotify Library")
	for i := 0; i < len(ids); i += 50 {
		upper := i + 49
		if upper >= len(ids) {
			upper = len(ids)
		}
		err := spotifyClient.AddTracksToLibrary(ids[i:upper]...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added song %d - %d to library \n", i, upper)
	}
}

// Functions to help debug by deleting entire Spotify Library
func getAllSongIDsSpotify(spotifyClient spotify.Client) []spotify.ID {
	songsRet := make([]spotify.ID, 0)

	songCount := getSongCountSpotify(spotifyClient)
	limit := 50
	for offset := 0; offset < songCount; offset += 50 {
		options := spotify.Options{Limit: &limit, Offset: &offset}
		songs, err := spotifyClient.CurrentUsersTracksOpt(&options)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range songs.Tracks {
			songsRet = append(songsRet, val.FullTrack.ID)
		}
	}

	return songsRet
}

func deleteAllSongsSpotify(spotifyClient spotify.Client) {
	idList := getAllSongIDsSpotify(spotifyClient)
	for i := 0; i < len(idList); i += 50 {
		if i+49 < len(idList) {
			spotifyClient.RemoveTracksFromLibrary(idList[i : i+49]...)
		} else {
			spotifyClient.RemoveTracksFromLibrary(idList[i:]...)
		}
	}

}
