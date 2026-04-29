package main

import (
	"log"
)

func main() {
	token, err := GetToken()
	if err != nil {
		log.Fatalln(err)
	}

	clips, err := GetClips(token)
	if err != nil {
		log.Fatalln(err)
	}

	DownloadClips(token, clips)

}
