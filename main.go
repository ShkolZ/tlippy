package main

import "log"

func main() {
	token, err := GetToken()
	if err != nil {
		log.Fatalln(err)
	}

}
