package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adlio/trello"
)

var mockTrelloServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/members/me/boards":
		jsonResponse, _ := json.Marshal(mockBoards)
		w.Write(jsonResponse)
	}
}))

func TestSelectBoard(t *testing.T) {
	client := testClient(mockTrelloServer.URL)

	out := &bytes.Buffer{}
	wantOut := `Select the Trello board you want to import from.
(0) Test
(1) Test Two
`
	selection := 0
	in := bytes.Buffer{}
	in.Write([]byte(fmt.Sprint(selection)))

	board := selectBoard(client, &in, out)

	assertStringMatch(board.ID, mockBoards[selection].ID, t)
	assertStringMatch(out.String(), wantOut, t)
}

var mockBoards = []trello.Board{
	{
		ID:   "1",
		Name: "Test",
	},
	{
		ID:   "1",
		Name: "Test Two",
	},
}

func testClient(baseURL string) *trello.Client {
	client := trello.NewClient(apikey, testtoken)
	client.BaseURL = baseURL
	return client
}
