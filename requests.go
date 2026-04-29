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
	Data []Game `json:"data"`
}

type Clip struct {
	ID          string `json:"id"`
	CreatorName string `json:"creator_name"`
	Title       string `json:"title"`
	Views       int    `json:"view_count"`
	CreatedAt   string `json:"created_at"`
}
type Clips struct {
	Clips []Clip `json:"data"`
}

func GetClips(token *Token) error {

	games := getGameId(token)

	query := url.Values{}
	query.Set("first", "50")
	query.Set("game_id", games.Data[0].Id)
	query.Set("started_at", time.Now().Add(-time.Hour*24*7).Format(time.RFC3339))

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/clips?%v", query.Encode())
	fmt.Println(endpoint)
	// query.Set("")
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	clips := Clips{}
	if err = json.Unmarshal(data, &clips); err != nil {
		fmt.Println(err)
	}
	fmt.Println(clips)
	return nil
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
