package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ShkolZ/tlippy/internal/helpers"
)

func DownloadClips(token *Token, clips Clips, cfg *Config) {
	for _, clip := range clips.Clips {
		split := strings.Split(clip.Url, "/")
		slug := split[len(split)-1]
		DownloadClip(clip, slug, cfg.DownloadPath)
	}
}

func DownloadClip(clip Clip, slug string, dPath string) error {
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

	clipName := fmt.Sprintf("[%v] %v-%v.mp4", helpers.FormatDate(clip.CreatedAt), helpers.CleanName(clip.CreatorName), helpers.CleanName(clip.Title))
	fullPath := path.Join(dPath, clipName)

	file, err := os.Create(fullPath)
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
