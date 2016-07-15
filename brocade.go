package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Login -> get Token
// List (rresource, token)

type BrocadeConfig struct {
	User     string
	Password string
	Token    string
	Endpoint string
}

func getConfig() *BrocadeConfig {
	user := os.Getenv("BROCADE_USER")
	password := os.Getenv("BROCADE_PASSWORD")
	endpoint := os.Getenv("BROCADE_ENDPOINT")
	token := os.Getenv("BROCADE_TOKEN")
	config := BrocadeConfig{
		User:     user,
		Password: password,
		Endpoint: endpoint,
		Token:    token,
	}

	return &config
}

func doLogin(config *BrocadeConfig, jsonBody string) (*BrocadeConfig, error) {

	path := fmt.Sprintf("%s/rest/login", config.Endpoint)
	var jsonStr = []byte(jsonBody)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonStr))

	req.Header.Set("WSUsername", config.User)
	req.Header.Set("WSPassword", config.Password)
	req.Header.Set("Accept", "application/vnd.brocade.networkadvisor+json;version=v1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	token := resp.Header.Get("WStoken")
	config.Token = token
	var e error
	if token == "" {
		e = errors.New("Authentication failed")
	}
	return config, e

}
