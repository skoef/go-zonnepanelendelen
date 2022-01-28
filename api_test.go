package zonnepanelendelen

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
