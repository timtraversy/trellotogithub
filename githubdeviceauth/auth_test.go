package githubdeviceauth

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

type MockClient struct{}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.String() {
	case "https://github.com/login/oauth/access_token":
		return &http.Response{
			StatusCode: 200,
		}, nil
	}
	return nil, nil
}

func (m *MockClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return nil, nil
}
func TestAuthenticateGithub(t *testing.T) {
	client := &MockClient{}
	want := github.Client{}
	got, err := AuthenticateGithub(client)
	if err != nil {
		t.Error("Unexpected error")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}
