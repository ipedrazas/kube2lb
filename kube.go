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
	Ingress   = "%s/apis/extensions/v1beta1/ingresses"
	Nodes     = "%s/api/v1/nodes"
	Services  = "%s/api/v1/services"
	Pods      = "%s/api/v1/pods"
	EndPoints = "%s/api/v1/endpoints"
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

func getItems(config Config, path string) (*ItemList, error) {
	var nodeList ItemList
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

func getEndPoints(config Config) (*ItemList, error) {

	path := fmt.Sprintf(EndPoints, config.ApiServer)
	nodeList, err := getItems(config, path)
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}

func getNodes(config Config) ([]string, error) {
	var nodes []string
	path := fmt.Sprintf(Nodes, config.ApiServer)
	nodeList, err := getItems(config, path)
	if err != nil {
		return nil, err
	}
	for _, item := range nodeList.Items {
		nodes = append(nodes, item.Metadata.Name)
	}
	return nodes, nil
}

func getPods(config Config) (*ItemList, error) {
	path := fmt.Sprintf(Pods, config.ApiServer)
	nodeList, err := getItems(config, path)
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}

// func getUnschedulable(config Config) (*ItemList, error) {
// 	var unscheduled []Node
// 	nodes, err := getNodes(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, item := range nodes.Items {
// 		if !item.Spec.Unschedulable {
// 			unscheduled = append(unscheduled, item)
// 		}
// 	}
// 	nodes.Items = unscheduled
// 	return nodes, nil
// }

func getPorts(config Config) ([]int, error) {
	var nodeList ItemList
	var exposedPorts []int
	path := fmt.Sprintf(Services, config.ApiServer)
	body, error := doGet(config, path)
	defer body.Close()
	if error != nil {
		return nil, error
	}
	err := json.NewDecoder(body).Decode(&nodeList)
	if err != nil {
		return nil, err
	}
	for _, item := range nodeList.Items {
		if item.Spec.Type == "NodePort" {
			for _, e := range item.Spec.Ports {
				exposedPorts = append(exposedPorts, e.NodePort)
			}
		}
	}
	return exposedPorts, nil
}

func getIngresses(config Config) (*ItemList, error) {
	path := fmt.Sprintf(Ingress, config.ApiServer)
	nodeList, err := getItems(config, path)
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}
