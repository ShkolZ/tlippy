package main

import (
	"strings"
)

func DownloadClips(token *Token, clips Clips) {
	// for _, clip := range clips.Clips {
	clip := clips.Clips[0]
	split := strings.Split(clip.Url, "/")
	slug := split[len(split)-1]

	DownloadClip(slug)
	// }
}

func DownloadClip(slug string) {
	GetClipLinks(slug)
}
