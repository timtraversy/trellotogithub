package main

import (
	"bufio"
	"fmt"
	"io"
)

// authenticateTrello auths Trello
func authenticateTrello(in io.Reader, out io.Writer) (token string, err error) {
	fmt.Println("Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Println("https://trello.com/1/authorize?expiration=never&name=Trello", "%", "20to", "%", "20Github&scope=read,account&response_type=token&key=", apikey)
	fmt.Println("When the authentication process is complete, enter the token you receied here:")
	scanner := bufio.NewReader(in)
	return scanner.ReadString('\n')
}
