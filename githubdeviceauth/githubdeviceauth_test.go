package githubdeviceauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/github"
)

var mockDeviceCodesResponse = deviceCodesResponse{
	DeviceCode:      "1",
	UserCode:        "2",
	VerificationURI: "test.com",
	Interval:        5,
	ExpiresIn:       900,
}

var mockDeviceAuthorizationResponse = deviceAuthorizationResponse{
	AccessToken: "1",
}

func TestAuthenticateGithub(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case deviceCodeURL:
			jsonResponse, _ := json.Marshal(mockDeviceCodesResponse)
			w.Write(jsonResponse)
		case authTokenURL:
			jsonResponse, _ := json.Marshal(mockDeviceAuthorizationResponse)
			w.Write(jsonResponse)
		}
	}))
	defer server.Close()

	gotOut := bytes.NewBufferString("")

	want := mockDeviceAuthorizationResponse.AccessToken

	got, err := AuthenticateGithub([]github.Scope{github.ScopeAdminOrg}, gotOut, server.URL)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}

	var wantOut = []string{
		fmt.Sprintln(requestingCodes),
		fmt.Sprintf(codeEntryInstructions, mockDeviceCodesResponse.VerificationURI, mockDeviceCodesResponse.UserCode),
		fmt.Sprintln(waitingForAuthorization),
	}

	for _, want := range wantOut {
		got, _ := gotOut.ReadString('\n')
		if got != want {
			t.Errorf("Got %v, want %v", gotOut, wantOut)
		}
	}
}
