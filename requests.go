package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Game struct {
	Id   string
	Name string
}

type Games struct {
	Data []Game
}

func GetClips(token *Token) {

	games := getGameId(token)

	query := url.Values{}
	query.Set("game_id", games.Data[0].Id)
	fmt.Println(time.Now())
	fmt.Println(time.Now().Add(-time.Hour * 24 * 7))
	// query.Set("")
	http.Get("https://api.twitch.tv/helix/clips")
}

func getGameId(token *Token) *Games {
	query := url.Values{}
	query.Set("igdb_id", "301298")

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/games?%v", query.Encode())

	fmt.Println(endpoint)
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	fmt.Println(string(data))
	games := Games{}
	err = json.Unmarshal(data, &games)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(games)

	return &games
}
