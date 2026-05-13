package oauth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var (
	clientID     = ""
	clientSecret = ""
)

type Token struct {
	Token    string `json:"access_token"`
	Expires  int    `json:"expires_in"`
	Type     string `json:"token_type"`
	ClientID string
}

func GetToken() (*Token, error) {

	if clientID == "" || clientSecret == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
		clientID = os.Getenv("CLIENT_ID")
		clientSecret = os.Getenv("CLIENT_SECRET")
	}

	query := url.Values{}
	query.Set("client_id", clientID)
	query.Set("client_secret", clientSecret)
	query.Set("grant_type", "client_credentials")

	res, err := http.PostForm("https://id.twitch.tv/oauth2/token", query)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}
	token.ClientID = clientID

	return &token, nil
}
