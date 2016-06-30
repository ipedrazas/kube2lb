package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	Ingress = "%s/apis/extensions/v1beta1/ingresses"
)

func getIngresses(config Config) (*NodeList, error) {
	var nodeList NodeList
	client := &http.Client{}
	if config.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	path := fmt.Sprintf(Ingress, config.ApiServer)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := fmt.Sprintf("Bearer %s", config.BearerToken)
	req.Header.Set("Authorization", bearer)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&nodeList)
	if err != nil {
		return nil, err
	}
	return &nodeList, nil
}
