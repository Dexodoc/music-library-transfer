package main

import (
	"math/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/zmb3/spotify"
)

type SpotifyClientBuilderConfig struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	LocalPort    string
}

// SpotifyClientBuilder builds an authenticated Spotify client
type SpotifyClientBuilder struct {
	Config *SpotifyClientBuilderConfig
	auth   spotify.Authenticator
	state  string
	ch     chan *spotify.Client
}

func NewSpotifyClientBuilder(config *SpotifyClientBuilderConfig) *SpotifyClientBuilder {
	c := &SpotifyClientBuilder{}
	if config == nil {
		c.Config = &SpotifyClientBuilderConfig{
			Scopes:       []string{spotify.ScopeUserTopRead},
			LocalPort:    "8080",
			ClientID:     os.Getenv("SPOTIFY_ID"),
			ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		}
	} else {
		c.Config = config
	}
	c.auth = spotify.NewAuthenticator(fmt.Sprintf("http://localhost:8080/callback"), c.Config.Scopes...)
	if c.Config.ClientID != "" && c.Config.ClientSecret != "" {
		c.auth.SetAuthInfo(c.Config.ClientID, c.Config.ClientSecret)
	}
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

func runSpotify() {
	spotifyClientBuilder := NewSpotifyClientBuilder(nil)
	spotifyClient, err := spotifyClientBuilder.GetClient()
	
	if err != nil {
		log.Fatal(err)
	}
	// get the logged in user
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello " + user.DisplayName)
	// get the user's top tracks
	tracks, err := spotifyClient.CurrentUsersTopTracks()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(tracks.Tracks); i++ {
		fmt.Println(tracks.Tracks[i].Name)
	}
	
    
}

