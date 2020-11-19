package main

import (
	"flag"
)

func main() {
	authCmd := flag.NewFlagSet("auth", flag.ExitOnError)
	authEnable := authCmd.Bool("enable", false, "description")
	// client := http.DefaultClient
	// _, err := AuthenticateTrello()
	// if err != nil {
	// 	return
	// }
	// cli, _ := AuthenticateGithub(client)
}
