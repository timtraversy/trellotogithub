package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/adlio/trello"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// AuthenticateTrello auths Trello
func AuthenticateTrello() (*trello.Client, error) {

	memberFromFile, err := readMember()
	if err == nil {
		return trello.NewClient(TrelloAPIKey, memberFromFile.Token), nil
	}

	fmt.Println("Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Println("https://trello.com/1/authorize?expiration=never&name=Trello", "%", "20to", "%", "20Github&scope=read,account&response_type=token&key=", TrelloAPIKey)
	fmt.Println("When the authentication process is complete, enter the token you receied here:")

	token := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	token <- scanner.Text()

	client := trello.NewClient(TrelloAPIKey, <-token)
	member, err := client.GetMember("me", trello.Defaults())

	if err != nil {
		fmt.Println("Failed to authenticate. Please check your token and try again")
		return nil, err
	}

	memberFromFile = storedMember{Token: client.Token, Name: member.Username}
	memberFromFile.writeMember()

	fmt.Println("Successfull authenticated as", member.Username)
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

// AuthenticateGithub auths with github
func AuthenticateGithub(client *http.Client) (*github.Client, error) {
	fmt.Println("Communicating with GitHub...")
	deviceAuthorization := getDeviceAuthorization(client)
	fmt.Printf("Please go to %s and enter the code %s\n", deviceAuthorization.VerificationURI, deviceAuthorization.UserCode)
	token := waitForAuthorization(client, deviceAuthorization.Interval, deviceAuthorization.ExpiresIn, deviceAuthorization.DeviceCode)
	fmt.Println("Token: ", token)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

func getDeviceAuthorization(client *http.Client) deviceAuthorizationCodeResponse {
	request := deviceAuthorizationCodeRequest{
		ClientID: GithubClientId,
		Scope:    "repo admin:org read:user",
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "https://github.com/login/device/code", bytes.NewBuffer(jsonBody))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var deviceAuthorization deviceAuthorizationCodeResponse
	json.Unmarshal(body, &deviceAuthorization)
	return deviceAuthorization
}

type deviceAuthorizationCodeRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

type deviceAuthorizationCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

func waitForAuthorization(client *http.Client, interval int, expiresIn int, deviceCode string) string {
	fmt.Println("Waiting for authorization in the browser...")
	for elapsed := 0; elapsed < expiresIn; elapsed += interval {
		request := deviceAuthorizationRequest{
			ClientID:   GithubClientId,
			DeviceCode: deviceCode,
			GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
		}
		jsonBody, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonBody))
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var deviceAuthorization deviceAuthorizationResponse
		json.Unmarshal(body, &deviceAuthorization)
		if (deviceAuthorization == deviceAuthorizationResponse{}) {
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}
		return deviceAuthorization.AccessToken
	}
	return ""
}

type deviceAuthorizationRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type deviceAuthorizationResponse struct {
	AccessToken string `json:"access_token"`
}
