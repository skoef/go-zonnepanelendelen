package zonnepanelendelen

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/skoef/go-zonnepanelendelen/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAPI(t *testing.T) {
	a := New("foo", "bar")

	assert.Equal(t, a.username, "foo")
	assert.Equal(t, a.password, "bar")
	assert.Nil(t, a.token)
}

func TestGetAPIURL(t *testing.T) {
	assert.Equal(t, "https://mijnstroom.zonnepanelendelen.nl/api/v1/foo?foo=bar", getAPIURL("/foo", "?foo=bar"))
	assert.Equal(t, "https://mijnstroom.zonnepanelendelen.nl/api/v1/foo?foo=bar", getAPIURL("foo", "foo=bar"))
}

func TestIsLoggedIn(t *testing.T) {
	a := API{}
	assert.False(t, a.isLoggedIn())

	a.token = &AuthToken{}
	assert.False(t, a.isLoggedIn())

	a.token.Token = "foo"
	assert.True(t, a.isLoggedIn())
}

func TestCall(t *testing.T) {
	t.Run("login first", func(t *testing.T) {
		// prepare API, without token
		// this should trigger login when other path is called
		api := API{}
		client := mocks.HTTPClient{}
		// login call
		tokenData, _ := os.ReadFile("testdata/token.json")
		client.On("Do", mock.IsType(&http.Request{})).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(tokenData)),
		}, nil).Once()
		// throw error on second call
		apiError := errors.New("test")
		client.On("Do", mock.IsType(&http.Request{})).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("")),
		}, apiError).Once()

		// set mock client as API client
		// and later revert this for further testing
		HTTPAPIClient = &client
		defer func() {
			HTTPAPIClient = nil
		}()

		// the actual call should fail, but we should be authenticated by now
		data, err := api.call("GET", "/path", "")
		assert.Equal(t, apiError, err)
		assert.Nil(t, data)
		if assert.NotNil(t, api.token) {
			assert.Equal(t, "John Doe", api.token.Name)
			assert.Equal(t, "1f0a796fd8342ca8d301c5f9a9a71e56", api.token.Token)
		}

		client.AssertExpectations(t)
	})
}

func TestGetProjects(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// assume logged in API client
		// prepare API, without token
		// this should trigger login when other path is called
		api := API{
			token: &AuthToken{Token: "test"},
		}
		client := mocks.HTTPClient{}
		// project return
		projectData, _ := os.ReadFile("testdata/projects.json")
		client.On("Do", mock.IsType(&http.Request{})).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(projectData)),
		}, nil).Once()

		// set mock client as API client
		// and later revert this for further testing
		HTTPAPIClient = &client
		defer func() {
			HTTPAPIClient = nil
		}()

		projects, err := api.GetProjects()
		assert.NoError(t, err)
		if assert.Equal(t, 1, len(projects)) {
			assert.Equal(t, "Test project", projects[0].Name)
			assert.Equal(t, 123, projects[0].ID)
			assert.Equal(t, 52.34899, projects[0].Latitude)
			assert.Equal(t, 4.92096, projects[0].Longitude)
			assert.Equal(t, true, projects[0].IsHidden)
		}

		client.AssertExpectations(t)
	})
}
