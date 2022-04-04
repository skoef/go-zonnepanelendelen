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

func TestGetProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// assume logged in API client
		// prepare API, without token
		// this should trigger login when other path is called
		api := API{
			token: &AuthToken{Token: "test"},
		}
		client := mocks.HTTPClient{}
		// project return
		projectData, _ := os.ReadFile("testdata/project.json")
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

		project, err := api.GetProject(123)
		assert.NoError(t, err)
		assert.Equal(t, "Test project", project.Name)
		assert.Equal(t, 123, project.ID)
		assert.Equal(t, 52.34899, project.Latitude)
		assert.Equal(t, 4.92096, project.Longitude)
		assert.Equal(t, true, project.IsHidden)
		assert.Equal(t, 123, project.Metrics.Parts)
		assert.Equal(t, 0.0246, project.Metrics.Interest)
		assert.Equal(t, 3075, project.Metrics.Value)

		if assert.Equal(t, 1, len(project.Metrics.Today.Measurements)) {
			assert.Equal(t, 0.123, project.Metrics.Today.Measurements[0].Production)
		}
		if assert.Equal(t, 31, len(project.Metrics.LastMonth.Measurements)) {
			assert.Equal(t, 1.976, project.Metrics.LastMonth.Measurements[0].Expected)
		}
		if assert.Equal(t, 12, len(project.Metrics.LastYear.Measurements)) {
			assert.Equal(t, 30.902, project.Metrics.LastYear.Measurements[2].Production)
		}
		assert.Equal(t, 0.01, project.Metrics.All.ROI)
		assert.EqualValues(t, 495, project.Metrics.All.TotalExpected)
		assert.Equal(t, 46.733, project.Metrics.All.TotalPower)
		assert.Equal(t, 3.422, project.Metrics.All.TotalProfit)

		if assert.Equal(t, 12, len(project.Metrics.All.Measurements)) {
			assert.Equal(t, 15.052, project.Metrics.All.Measurements[1].Cumulative)
			assert.EqualValues(t, 25.529, project.Metrics.All.Measurements[1].Expected)
			assert.EqualValues(t, 25.529, project.Metrics.All.Measurements[1].ExpectedCumulative)
			assert.Equal(t, "2022-02-01T02:00:00Z", project.Metrics.All.Measurements[1].Timestamp)
		}

		client.AssertExpectations(t)
	})
}
