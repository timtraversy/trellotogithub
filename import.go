package main

import (
	"fmt"
	"io"

	"github.com/adlio/trello"
)

func selectBoard(client *trello.Client, in io.Reader, out io.Writer) *trello.Board {
	boards, _ := client.GetMyBoards(trello.Defaults())
	fmt.Fprintln(out, "Select the Trello board you want to import from.")
	for i, board := range boards {
		fmt.Fprintf(out, "(%v) %v\n", i, board.Name)
	}
	var selection int
	fmt.Fscanln(in, selection)
	return boards[selection]
}
