package zonnepanelendelen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiHost      = "mijnstroom.zonnepanelendelen.nl"
	apiPrefix    = "/api/v1"
	apiLoginPath = "/obtain-auth-token/"
	authHeader   = "Authorization"
)

// HTTPAPIClient is the HTTP client used for interfacing with the API
var HTTPAPIClient HTTPClient

// New returns a new API struct for given username and password
func New(username, password string) API {
	return API{
		username: username,
		password: password,
	}
}

func getAPIURL(path string, query string) string {
	apiURL := url.URL{}
	apiURL.Scheme = "https"
	apiURL.Host = apiHost
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasPrefix(query, "?") {
		query = strings.Replace(query, "?", "", 1)
	}
	apiURL.RawQuery = query
	apiURL.Path = apiPrefix + path

	return apiURL.String()
}

func (a *API) call(method, path, data string) ([]byte, error) {
	// if API hasn't logged in yet, abort
	if !a.isLoggedIn() && path != apiLoginPath {
		if err := a.login(); err != nil {
			return nil, err
		}
	}

	// prepare HTTP request
	var url string
	var req *http.Request
	var err error

	switch method {
	case "POST":
		url = getAPIURL(path, "")
		req, err = http.NewRequest(method, url, strings.NewReader(data))
	case "GET":
		url = getAPIURL(path, data)
		req, err = http.NewRequest(method, url, nil)
	}

	// we couldn't prepare the request
	if err != nil {
		return nil, err
	}

	// make sure to set proper content-type when sending data
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// add Authorization header if we're not currently trying to login
	if path != apiLoginPath {
		req.Header[authHeader] = []string{a.token.String()}
	}

	if HTTPAPIClient == nil {
		HTTPAPIClient = http.DefaultClient
	}

	// do HTTP request
	resp, err := HTTPAPIClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200: // do nothing, this is OK
	default:
		return nil, fmt.Errorf("API HTTP error: %s", resp.Status)
	}

	// read and return entire response body
	return io.ReadAll(resp.Body)
}

func (a API) isLoggedIn() bool {
	if a.token == nil || a.token.Token == "" {
		return false
	}

	return true
}

func (a *API) login() error {
	// prepare form data
	params := url.Values{}
	params.Add("username", a.username)
	params.Add("password", a.password)

	// do API call
	data, err := a.call("POST", apiLoginPath, params.Encode())
	if err != nil {
		return err
	}

	// unmarshal JSON data and set token
	var token AuthToken
	err = json.Unmarshal(data, &token)
	if err != nil {
		return err
	}

	a.token = &token

	return nil
}

// GetProjects returns all projects that the authenticated account invested in
func (a API) GetProjects() ([]Project, error) {
	data, err := a.call("GET", "/projects/", "?view=index_only")
	if err != nil {
		return nil, err
	}

	var v struct {
		Projects []Project `json:"projects_invested_in"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v.Projects, nil
}

func (a API) GetProject(projectID int) (Project, error) {
	data, err := a.call("GET", fmt.Sprintf("project/%d", projectID), "")
	if err != nil {
		return Project{}, err
	}

	var v struct {
		Project Project `json:"project"`
		Metrics Metrics `json:"metrics"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return Project{}, err
	}

	v.Project.Metrics = v.Metrics

	return v.Project, nil
}
