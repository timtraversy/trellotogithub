package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adlio/trello"
)

func main() {
	client, err := authenticate()
	if err != nil {
		return
	}
	fmt.Println(client.Token)
}

func authenticate() (*trello.Client, error) {
	const APIKey = "4960d7908c097bc466ab6fb8074340d6"

	memberFromFile, err := readMember()
	if err == nil {
		return trello.NewClient(APIKey, memberFromFile.Token), nil
	}

	fmt.Println("Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Printf("https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%s\n", APIKey)
	fmt.Println("When the authentication process is complete, enter the token you receied here:")

	token := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	token <- scanner.Text()

	client := trello.NewClient(APIKey, <-token)
	member, err := client.GetMember("me", trello.Defaults())

	if err != nil {
		fmt.Println("Failed to authenticate. Please check your token and try again")
		return nil, err
	}

	memberFromFile = storedMember{Token: client.Token, Name: member.Username}
	memberFromFile.writeMember()

	fmt.Printf("Successfull authenticated as %s\n", member.Username)
	return client, nil
}

type storedMember struct {
	Token, Name string
}

const memberFileName = ".member.json"

func (s storedMember) writeMember() {
	file, _ := json.MarshalIndent(s, "", " ")
	ioutil.WriteFile(memberFileName, file, 0664)
}

func readMember() (storedMember, error) {
	bytes, err := ioutil.ReadFile(memberFileName)
	if err != nil {
		return storedMember{}, err
	}
	var data map[string]string
	json.Unmarshal(bytes, &data)
	return storedMember{Name: data["Name"], Token: data["Token"]}, nil
}
