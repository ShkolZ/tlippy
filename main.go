package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("not enough arguments")
	}

	cfg, err := SetConfig(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	token, err := GetToken()
	if err != nil {
		log.Fatalln(err)
	}

	clips, err := GetClips(token, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	DownloadClips(token, clips, cfg)

}
