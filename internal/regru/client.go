package regru

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const BASE_URL = "https://api.reg.ru/api/regru2"

func NewClient(username, password, certPem, keyPem string) (*Client, error) {
	var (
		nilClient *Client
		transport = http.Transport{}
		auth      = clientAuth{username, password, nil}
	)

	if certPem != "" && keyPem != "" {
		cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
		if err != nil {
			return nilClient, err
		}

		transport = http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}

		decodedPem, _ := pem.Decode([]byte(keyPem))
		unknownKey, err := x509.ParsePKCS8PrivateKey(decodedPem.Bytes)
		if err != nil {
			return nilClient, err
		}

		privateKey, ok := unknownKey.(*rsa.PrivateKey)
		if !ok {
			return nilClient, fmt.Errorf("bad private key")
		}

		auth.privateKey = privateKey
	}

	httpClient := http.Client{
		Timeout:   60 * time.Second,
		Transport: &transport,
	}

	return &Client{auth, httpClient}, nil
}

type clientAuth struct {
	username   string
	password   string
	privateKey *rsa.PrivateKey
}

type Client struct {
	auth       clientAuth
	httpClient http.Client
}

func (c *Client) fetch(action Action, data map[string]string) (Response, error) {
	res := Response{}

	uri, err := url.JoinPath(BASE_URL, string(action))
	if err != nil {
		return res, err
	}

	form := url.Values{}
	values := make([]string, 0, len(data)+1)

	for key, value := range data {
		form.Add(key, value)
		values = append(values, value)
	}

	values = append(values, c.auth.username)

	form.Add("username", c.auth.username)

	if c.auth.password != "" {
		form.Add("password", c.auth.password)
	} else {
		sort.Strings(values)

		hash := sha512.Sum512([]byte(strings.Join(values, ";")))

		sig, err := rsa.SignPKCS1v15(nil, c.auth.privateKey, crypto.SHA512, hash[:])
		if err != nil {
			return res, err
		}

		form.Add("sig", base64.StdEncoding.EncodeToString(sig))
	}

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
