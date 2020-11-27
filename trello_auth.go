package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func authenticateTrello(in io.Reader, out io.Writer) (token string, err error) {
	fmt.Fprintln(out, "Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Fprintf(out, "https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%v\n", apikey)
	fmt.Fprintln(out, "When the authentication process is complete, enter the token you receied here:")

	tokenChn := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	tokenChn <- scanner.Text()
	return <-tokenChn, nil
}
