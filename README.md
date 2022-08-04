# Music Library Transfer

A basic CLI tool to be able to transfer entire libraries of songs from Apple Music to Spotify and vice versa. Makes use of API's that both services offer to set all song data for a user and to transfer it to the other library. 

Authorisation happens in a browser window to generate user tokens for both services.

# Usage
You would need to go get your own developer token from apple and setup a private key to generate a JWT token and similary for Spotify you would need to setup a free app and get your cliet secret and client id and place them in main.go and also in index.html as needed.

Run with `go run .`