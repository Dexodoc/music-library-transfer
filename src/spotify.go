package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

 func spotify() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.music.apple.com/v1/me/library/songs", nil)

	// Setting headers for user and developer authentication
    	

    req.Header.Set("Music-User-Token", userToken)

	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	req.URL.RawQuery = q.Encode()

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &r)
	
       

}

