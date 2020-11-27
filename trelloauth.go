package main

import (
	"fmt"
	"io"
)

func authenticateTrello(in io.Reader, out io.Writer) (token string, err error) {
	fmt.Fprintln(out, "Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Fprintln(out, "https://trello.com/1/authorize?expiration=never&name=Trello", "%", "20to", "%", "20Github&scope=read,account&response_type=token&key=", apikey)
	fmt.Fprintln(out, "When the authentication process is complete, enter the token you receied here:")
	fmt.Fscanln(in, token)
	return
}
