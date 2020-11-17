package main

import (
	"fmt"
	"testing"

	"github.com/adlio/trello"
)

func TestGetBoard(t *testing.T) {
	client := trello.NewClient(APIKey, ServerKey)
	boards, _ := client.GetMyBoards(trello.Defaults())
	fmt.Println(boards[0].Name)
}
