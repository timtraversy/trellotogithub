package deviceauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-github/v32/github"
)

func TestAuthenticateGithub(t *testing.T) {
	gotOut := &bytes.Buffer{}

	var githubAuthenticator = GithubDeviceAuthenticator{
		BaseURL: mockGithubServer.URL,
		Out:     gotOut,
	}

	got, err := githubAuthenticator.AuthenticateGithub([]github.Scope{github.ScopeAdminOrg})
	checkError(err, t)
	compareStructs(got, mockDeviceAuthorizationResponse.AccessToken, t)

	var wantOut = []string{
		fmt.Sprintln(requestingCodes),
		fmt.Sprintf(codeEntryInstructions, mockDeviceCodes.VerificationURI, mockDeviceCodes.UserCode),
	}

	for _, want := range wantOut {
		got, _ := gotOut.ReadString('\n')
		if got != want {
			t.Errorf("Got %v, want %v", gotOut, wantOut)
		}
	}
}

func TestRequestDeviceCodes(t *testing.T) {
	got, err := githubAuthenticator.RequestDeviceCodes([]github.Scope{github.ScopeAdminOrg})
	checkError(err, t)
	compareStructs(got, mockDeviceCodes, t)
}

func TestWaitForAuthorization(t *testing.T) {
	got, err := githubAuthenticator.WaitForAuthorization(mockDeviceCodes)
	checkError(err, t)
	compareStructs(got, mockDeviceAuthorizationResponse.AccessToken, t)
}

func checkError(err error, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}
}

func compareStructs(got interface{}, want interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

var mockDeviceCodes = DeviceCodes{
	DeviceCode:      "1",
	UserCode:        "2",
	VerificationURI: "test.com",
	Interval:        5,
	ExpiresIn:       900,
}

var mockDeviceAuthorizationResponse = deviceAuthorizationResponse{
	AccessToken: "1",
}

var mockGithubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case deviceCodeURL:
		jsonResponse, _ := json.Marshal(mockDeviceCodes)
		w.Write(jsonResponse)
	case authTokenURL:
		jsonResponse, _ := json.Marshal(mockDeviceAuthorizationResponse)
		w.Write(jsonResponse)
	}
}))

var githubAuthenticator = GithubDeviceAuthenticator{
	BaseURL: mockGithubServer.URL,
}
