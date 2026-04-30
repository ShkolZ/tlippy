package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type ClipTokenResponse struct {
	Data struct {
		Clip struct {
			PlaybackAccessToken struct {
				Signature string `json:"signature"`
				Value     string `json:"value"`
			} `json:"playbackAccessToken"`

			VideoQualities []struct {
				Quality   string `json:"quality"`
				SourceURL string `json:"sourceURL"`
			} `json:"videoQualities"`
		} `json:"clip"`
	} `json:"data"`
}

type Game struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	Url         string `json:"url"`
	Thumbnail   string `json:"thumbnail_url"`
}
type Clips struct {
	Clips []Clip `json:"data"`
}

const gqlURI = "https://gql.twitch.tv/gql"
const clientID = "kimne78kx3ncx6brgo4mv6wki5h1ko"

func GetClips(token *Token, cfg *Config) (Clips, error) {
	games := getGameId(token)

	query := url.Values{}
	query.Set("first", strconv.Itoa(cfg.ClipsAmount))
	query.Set("game_id", games.Data[0].Id)
	query.Set("started_at", time.Now().Add(-time.Hour*24*7).Format(time.RFC3339))

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/clips?%v", query.Encode())

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

	return clips, nil
}

func getGameId(token *Token) *Games {
	query := url.Values{}
	query.Set("igdb_id", "301298")

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/games?%v", query.Encode())

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	games := Games{}
	err = json.Unmarshal(data, &games)
	if err != nil {
		fmt.Println(err)
	}

	return &games
}

func GetClipLinks(slug string) (*ClipTokenResponse, error) {
	// hash := sha256.New()
	body := map[string]any{
		"operationName": "VideoAccessToken_Clip",
		"variables": map[string]any{
			"slug": slug,
		},
		"extensions": map[string]any{
			"persistedQuery": map[string]any{
				"version":    1,
				"sha256Hash": "36b89d2507fce29e5ca551df756d27c1cfe079e2609642b4390aa4c35796eb11",
			},
		},
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", gqlURI, bytes.NewReader(b))
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(data))
	ctr := ClipTokenResponse{}
	json.Unmarshal(data, &ctr)

	return &ctr, nil

}
