package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	Ingress = "%s/apis/extensions/v1beta1/ingresses"
	Nodes   = "%s/api/v1/nodes"
)

func doGet(config Config, path string) (io.ReadCloser, error) {

	client := &http.Client{}
	if config.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

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
	return resp.Body, nil

}

func getNodes(config Config) (*ItemList, error) {
	var nodeList ItemList
	path := fmt.Sprintf(Nodes, config.ApiServer)
	body, error := doGet(config, path)
	defer body.Close()
	if error != nil {
		return nil, error
	}
	err := json.NewDecoder(body).Decode(&nodeList)
	if err != nil {
		return nil, err
	}
	return &nodeList, nil
}

func getIngresses(config Config) (*ItemList, error) {
	var nodeList ItemList
	path := fmt.Sprintf(Ingress, config.ApiServer)
	body, error := doGet(config, path)
	defer body.Close()
	if error != nil {
		return nil, error
	}
	err := json.NewDecoder(body).Decode(&nodeList)
	if err != nil {
		return nil, err
	}
	return &nodeList, nil
}
