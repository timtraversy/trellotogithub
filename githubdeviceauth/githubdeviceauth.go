package githubdeviceauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
)

type GithubDeviceAuthenticator struct {
	BaseURL string
	Out     io.Writer
}

func NewAuthenticator() *GithubDeviceAuthenticator {
	return &GithubDeviceAuthenticator{
		BaseURL: "https://github.com",
		Out:     os.Stdout,
	}
}

const requestingCodes = "Requesting device and user verification codes from GitHub"
const codeEntryInstructions = "Please go to %s and enter the code %s\n"

func (g *GithubDeviceAuthenticator) AuthenticateGithub(scopes []github.Scope) (token string, err error) {
	fmt.Fprintln(g.Out, requestingCodes)
	deviceCodes, err := g.RequestDeviceCodes(scopes)
	fmt.Fprintf(g.Out, codeEntryInstructions, deviceCodes.VerificationURI, deviceCodes.UserCode)
	token, err = g.WaitForAuthorization(deviceCodes)
	return
}

const deviceCodeURL = "/login/device/code"

type deviceCodesRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

type DeviceCodes struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

func (g *GithubDeviceAuthenticator) RequestDeviceCodes(scopes []github.Scope) (deviceCodes DeviceCodes, err error) {
	var scopeString string
	for _, scope := range scopes {
		scopeString += fmt.Sprintf("%s ", scope)
	}
	request := deviceCodesRequest{
		ClientID: clientID,
		Scope:    scopeString,
	}
	g.postRequest(request, deviceCodeURL, &deviceCodes)
	return
}

const authTokenURL = "/login/oauth/access_token"

type deviceAuthorizationRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type deviceAuthorizationResponse struct {
	AccessToken string `json:"access_token"`
}

func (g *GithubDeviceAuthenticator) WaitForAuthorization(deviceCodes DeviceCodes) (token string, err error) {
	for elapsed := 0; elapsed < deviceCodes.ExpiresIn; elapsed += deviceCodes.Interval {
		request := deviceAuthorizationRequest{
			ClientID:   clientID,
			DeviceCode: deviceCodes.DeviceCode,
			GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
		}
		var deviceAuthorization deviceAuthorizationResponse
		g.postRequest(request, authTokenURL, &deviceAuthorization)
		if (deviceAuthorization == deviceAuthorizationResponse{}) {
			time.Sleep(time.Duration(deviceCodes.Interval) * time.Second)
			continue
		}
		return deviceAuthorization.AccessToken, nil
	}
	return "", errors.New("Authorization timed out")
}

func (g *GithubDeviceAuthenticator) postRequest(requestBody interface{}, url string, resultStructPtr interface{}) error {
	jsonBody, _ := json.Marshal(requestBody)
	request, _ := http.NewRequest("POST", g.BaseURL+url, bytes.NewBuffer(jsonBody))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, resultStructPtr)
	return nil
}
