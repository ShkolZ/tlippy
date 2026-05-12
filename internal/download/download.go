package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/helpers"
	"github.com/ShkolZ/tlippy/internal/oauth"
	"github.com/ShkolZ/tlippy/internal/request"
)

type Progress struct {
	Current int
	Total   int
	Name    string
	Err     error
	Done    bool
}

func StartDownloadChan(input *config.UserInput) <-chan Progress {
	ch := make(chan Progress, 1)

	go func() {
		defer close(ch)

		token, err := oauth.GetToken()
		if err != nil {
			ch <- Progress{Err: err, Done: true}
			return
		}

		if input.Mode == config.ModeBulk {
			clips, err := request.GetClips(token, input)
			if err != nil {
				ch <- Progress{Err: err, Done: true}
				return
			}
			total := len(clips.Clips)
			for i, clip := range clips.Clips {
				split := strings.Split(clip.Url, "/")
				slug := split[len(split)-1]
				dlErr := downloadClip(clip, slug, input.DownloadPath)
				ch <- Progress{
					Current: i + 1,
					Total:   total,
					Name:    clip.Title,
					Err:     dlErr,
					Done:    i+1 == total,
				}
			}
		} else {
			clip, err := request.GetClip(token, input.ClipID)
			if err != nil {
				ch <- Progress{Err: err, Done: true}
				return
			}
			dlErr := downloadClip(clip, input.ClipID, input.DownloadPath)
			ch <- Progress{
				Current: 1,
				Total:   1,
				Name:    clip.Title,
				Err:     dlErr,
				Done:    true,
			}
		}
	}()

	return ch
}

func downloadClip(clip request.Clip, slug string, dPath string) error {
	ctr, err := request.GetClipLinks(slug)
	if err != nil {
		return err
	}
	baseUrl := ctr.Data.Clip.VideoQualities[0].SourceURL
	downLink := fmt.Sprintf("%v?sig=%v&token=%v", baseUrl, ctr.Data.Clip.PlaybackAccessToken.Signature, url.QueryEscape(ctr.Data.Clip.PlaybackAccessToken.Value))
	resp, err := http.Get(downLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	clipName := fmt.Sprintf("[%v] %v-%v.mp4", helpers.FormatDate(clip.CreatedAt), helpers.CleanName(clip.CreatorName), helpers.CleanName(clip.Title))
	fullPath := path.Join(dPath, clipName)

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
