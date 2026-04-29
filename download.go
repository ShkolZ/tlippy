package main

func DownloadClips(token *Token, clips Clips) {
	for _, clip := range clips.Clips {
		DownloadClip(token, clip)
	}
}

func DownloadClip(token *Token, clip Clip) {
	DownloadRequest(token, clip)
}
