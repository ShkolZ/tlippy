package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DownloadClips(token *Token, clips Clips) {
	for _, clip := range clips.Clips {
		split := strings.Split(clip.Url, "/")
		slug := split[len(split)-1]
		DownloadClip(slug)
	}
}

func DownloadClip(slug string) error {
	ctr, err := GetClipLinks(slug)
	if err != nil {
		return err
	}
	baseUrl := ctr.Data.Clip.VideoQualities[0].SourceURL
	downLink := fmt.Sprintf("%v?sig=%v&token=%v", baseUrl, ctr.Data.Clip.PlaybackAccessToken.Signature, url.QueryEscape(ctr.Data.Clip.PlaybackAccessToken.Value))
	resp, err := http.Get(downLink)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	file, err := os.Create("zalupa")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
