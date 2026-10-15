package regru

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const BASE_URL = "https://api.reg.ru/api/regru2"

type Client struct {
	auth        ClientAuth
	httpClient  http.Client
	cachedZones map[string]*Zone
}

type ClientAuth struct {
	username string
	password string
}

func NewClient(username, password string) Client {
	auth := ClientAuth{username, password}

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	cachedZones := map[string]*Zone{}

	return Client{auth, httpClient, cachedZones}
}

func (c *Client) fetch(action string, data map[string]string) ([]byte, error) {
	uri, err := url.JoinPath(BASE_URL, action)
	if err != nil {
		return nil, err
	}

	form := url.Values{}

	for key, value := range data {
		form.Add(key, value)
	}

	form.Add("username", c.auth.username)
	form.Add("password", c.auth.password)

	response, err := c.httpClient.PostForm(uri, form)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (c *Client) setCachedZone(key string, zone Zone) {
	c.cachedZones[key] = &zone
}

func (c *Client) deleteCachedZone(key string) {
	c.cachedZones[key] = nil
}
