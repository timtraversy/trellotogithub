package githubdeviceauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// AuthenticateGithub auths with github
func AuthenticateGithub(client HTTPClient) (*github.Client, error) {
	fmt.Println("Communicating with GitHub...")
	deviceAuthorization := getDeviceAuthorization(client)
	fmt.Printf("Please go to %s and enter the code %s\n", deviceAuthorization.VerificationURI, deviceAuthorization.UserCode)
	token := waitForAuthorization(client, deviceAuthorization.Interval, deviceAuthorization.ExpiresIn, deviceAuthorization.DeviceCode)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

func getDeviceAuthorization(client HTTPClient) deviceAuthorizationCodeResponse {
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

func waitForAuthorization(client HTTPClient, interval int, expiresIn int, deviceCode string) string {
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
