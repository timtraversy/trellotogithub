package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/timtraversy/trello-to-github/githubdeviceauth"
)

func main() {
	// authCmd := flag.NewFlagSet("auth", flag.ExitOnError)
	// authEnable := authCmd.Bool("enable", false, "description")
	client := http.DefaultClient
	output := os.Stdout
	// _, err := AuthenticateTrello()
	// if err != nil {
	// 	return
	// }
	scopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
	token, err := githubdeviceauth.AuthenticateGithub(client, scopes, output)
	fmt.Println(token, err)
}
