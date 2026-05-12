package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/oauth"
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

func GetClip(token *oauth.Token, clipID string) (Clip, error) {
	query := url.Values{}
	query.Set("id", clipID)

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/clips?%v", query.Encode())

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Clip{}, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Clip{}, fmt.Errorf("reading response: %w", err)
	}

	clips := Clips{}
	if err = json.Unmarshal(data, &clips); err != nil {
		return Clip{}, fmt.Errorf("parsing response: %w", err)
	}

	if len(clips.Clips) == 0 {
		return Clip{}, fmt.Errorf("clip %q not found", clipID)
	}

	return clips.Clips[0], nil
}

func GetClips(token *oauth.Token, input *config.UserInput) (Clips, error) {
	games := getGameId(token, input.QueryName)

	query := url.Values{}

	query.Set("first", input.ClipCount)
	query.Set("game_id", games.Data[0].Id)
	query.Set("started_at", time.Now().Add(-getTime(input.TimeRange)).Format(time.RFC3339))

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/clips?%v", query.Encode())

	// query.Set("")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return Clips{}, err
	}
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

func getGameId(token *oauth.Token, name string) *Games {
	query := url.Values{}
	query.Set("name", name)

	endpoint := fmt.Sprintf("https://api.twitch.tv/helix/games?%v", query.Encode())

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}
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

func getTime(timeRange config.TimeRange) time.Duration {
	switch timeRange {
	case config.TimeRange24h:
		return time.Hour * 24
	case config.TimeRange7d:
		return time.Hour * 24 * 7
	case config.TimeRangeAll:
		return time.Hour * 24 * 365 * 100 // approx
	}
	return 0
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
