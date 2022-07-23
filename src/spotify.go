package main

import (
	"math/rand"
	"fmt"
	"log"
	"net/http"
	"github.com/zmb3/spotify"
)

type SpotifyClientBuilder struct {
	auth   spotify.Authenticator
	state  string
	ch     chan *spotify.Client
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
	client := <- c.ch
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

func runSpotify() {
	spotifyClient, err := NewSpotifyClientBuilder().GetClient()
	
	if err != nil {
		log.Fatal(err)
	}

	

	// get the logged in user
	user, err := spotifyClient.CurrentUser()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello " + user.DisplayName)
	k := 50
	o := spotify.Options{Limit: &k}

	tracks, err := spotifyClient.CurrentUsersTracksOpt(&o)

	if err != nil {
		log.Fatal(err)
	}

	for i, v := range tracks.Tracks{
		fmt.Println(i, v)
	}

	fmt.Println(tracks.Total)
}

