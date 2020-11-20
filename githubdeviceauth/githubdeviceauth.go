package githubdeviceauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

var githubBaseURL = "https://github.com"

func AuthenticateGithub(scopes []github.Scope, out io.Writer, baseURL string) (token string, err error) {
	if baseURL != "" {
		githubBaseURL = baseURL
	}
	deviceCodes := requestDeviceCodes(scopes, out)
	token = waitForAuthorization(deviceCodes.Interval, deviceCodes.ExpiresIn, deviceCodes.DeviceCode, out)
	return
}

const deviceCodeURL = "/login/device/code"
const requestingCodes = "Requesting device and user verification codes from GitHub"
const codeEntryInstructions = "Please go to %s and enter the code %s\n"

func requestDeviceCodes(scopes []github.Scope, out io.Writer) (deviceCodes deviceCodesResponse) {
	fmt.Fprintln(out, requestingCodes)
	var scopeString string
	for _, scope := range scopes {
		scopeString += fmt.Sprintf("%s ", scope)
	}
	request := deviceCodesRequest{
		ClientID: GithubClientId,
		Scope:    scopeString,
	}
	postRequest(request, deviceCodeURL, &deviceCodes)
	fmt.Fprintf(out, codeEntryInstructions, deviceCodes.VerificationURI, deviceCodes.UserCode)
	return
}

type deviceCodesRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

type deviceCodesResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

const authTokenURL = "/login/oauth/access_token"
const waitingForAuthorization = "Waiting for authorization in the browser..."

func waitForAuthorization(interval int, expiresIn int, deviceCode string, out io.Writer) string {
	fmt.Fprintln(out, waitingForAuthorization)
	for elapsed := 0; elapsed < expiresIn; elapsed += interval {
		request := deviceAuthorizationRequest{
			ClientID:   GithubClientId,
			DeviceCode: deviceCode,
			GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
		}
		var deviceAuthorization deviceAuthorizationResponse
		postRequest(request, authTokenURL, &deviceAuthorization)
		if (deviceAuthorization == deviceAuthorizationResponse{}) {
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}
		return deviceAuthorization.AccessToken
	}
	return ""
}

func postRequest(requestBody interface{}, url string, resultStructPtr interface{}) []byte {
	jsonBody, _ := json.Marshal(requestBody)
	request, _ := http.NewRequest("POST", githubBaseURL+url, bytes.NewBuffer(jsonBody))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, resultStructPtr)
	return responseBody
}

type deviceAuthorizationRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type deviceAuthorizationResponse struct {
	AccessToken string `json:"access_token"`
}
