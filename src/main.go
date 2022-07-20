package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
    Next string
    Data []Song
    Meta Meta
}

type Meta struct{
    Total int
}

type Song struct {
    Attributes Attributes
}

type Attributes struct {
    ArtistName string
    Name string
}

func main() {
	
    var appleDeveloperToken = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlRSVUM3NzdEWkwifQ.eyJpYXQiOjE2NTgzMzYzNzksImV4cCI6MTY3Mzg4ODM3OSwiaXNzIjoiM1M1OVVGV1A4WSJ9.ywmXqlzN700OszGHLXePHOk2ifJ3EFWk_mF5XZT8w9bPgK9epenRO62tbpjdVMNgu6p6Hio6095Wig73_3xtwA";
    
    // Currently personal token
    var musicUserToken = "AmIKSqcLaVo8ziS52+UVu6Qf+OW8vg8qXAhbpRi/+eh0ObM4AZ22pTVIZdxhsbpHCyiIOtjbmA8jKRkEqQL3+yo9AQSAJgTzu9l1CjSY1WLPOvl71B5bi9ZpFE/FhKahuqyP9Lf+aHV9H3bhThaPrWTrLO1HjxMbaxWjTzgoJVzGqdbviQdfYzxcHTe4nAzBNxKHysAm00SGl8xaVviEzb4Gs5sN1cpRrbso2TjO/cvlfvX+9g==";

    

/*    client := &http.Client{}*/
    /*req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/tracks?market=ES", nil)*/
    /*req.Header.Set("Authorization", token)*/
    /*res, _ := client.Do(req)*/
    /*body, _ := ioutil.ReadAll(res.Body)*/
    /*fmt.Print(string(body))*/

    client := &http.Client{}
    req, _ := http.NewRequest("GET", "https://api.music.apple.com/v1/me/library/songs", nil)
    
    // Setting headers for user and developer authentication
    req.Header.Set("Authorization", "Bearer " + appleDeveloperToken)
    req.Header.Set("Music-User-Token", musicUserToken)
    
    // Setting query parameters
    q := req.URL.Query()
    q.Add("limit", "100")
    q.Add("offset", "99")
    req.URL.RawQuery = q.Encode()

    res, _ := client.Do(req)
    resData := Response{}
    body, _ := ioutil.ReadAll(res.Body)
    json.Unmarshal(body, &resData)
    //fmt.Println(resData)
    
    fmt.Println(resData.Meta.Total);

    for _, v := range resData.Data {
        fmt.Println(v.Attributes.ArtistName)
        fmt.Println(v.Attributes.Name)

        fmt.Println()
    } 
}
