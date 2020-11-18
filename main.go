package main

import "net/http"

func main() {
	client := http.DefaultClient
	// _, err := AuthenticateTrello()
	// if err != nil {
	// 	return
	// }
	AuthenticateGithub(client)
}
