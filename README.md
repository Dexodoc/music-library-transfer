# Music Library Transfer

A basic CLI tool to be able to transfer entire libraries of songs from Apple Music to Spotify and vice versa. Makes use of API's that both services offer to set all song data for a user and to transfer it to the other library. 

Authorisation happens in a browser window to generate user tokens for both services.

# Usage

## Spotify
First, you would need to go to [Spotify developer portal](https://developer.spotify.com/dashboard) and create a new application and set the callback URI to be `http://localhost:8080/callback`. Then copy the *Client ID* and the *Client Secret* into `main.go`.

## Apple Music
You would need a developer account with apple and to generate a private key through their developer portal with access to music services. Once you have your private key `.p8` file you can place it in the jwt folder and rename it `AuthKey.p8` and then run `generate_jwt.js` filled in with a couple more details accessible from the developer portal to then generate a JWT to be used in `main.go` and `index.html`. 

## Run
`go run .`