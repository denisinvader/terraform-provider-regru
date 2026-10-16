package regru

import (
	"encoding/json"
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

func (c *Client) fetch(action Action, data map[string]string) (Response, error) {
	res := Response{}

	uri, err := url.JoinPath(BASE_URL, string(action))
	if err != nil {
		return res, err
	}

	form := url.Values{}

	for key, value := range data {
		form.Add(key, value)
	}

	form.Add("username", c.auth.username)
	form.Add("password", c.auth.password)

	response, err := c.httpClient.PostForm(uri, form)
	if err != nil {
		return res, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return res, err
	}

	if res.Result == ResultError {
		return res, RegruError{res.ErrorCode, res.ErrorText}
	}

	return res, nil
}

func (c *Client) setCachedZone(key string, zone Zone) {
	c.cachedZones[key] = &zone
}

func (c *Client) deleteCachedZone(key string) {
	c.cachedZones[key] = nil
}
