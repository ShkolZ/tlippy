package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Token struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
	Type    string `json:"token_type"`
}

func GetToken() (*Token, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	cid := os.Getenv("CLIENT_ID")
	cs := os.Getenv("CLIENT_SECRET")
	fmt.Println(cid, cs)

	query := url.Values{}
	query.Set("client_id", cid)
	query.Set("client_secret", cs)
	query.Set("grant_type", "client_credentials")

	res, err := http.PostForm("https://id.twitch.tv/oauth2/token", query)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}
	fmt.Println(token)

	return &token, nil
}
