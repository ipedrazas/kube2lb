package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	fmt.Printf("Kube2lb, an utility to update load balancers\n")

	// read params from env vars
	apiserver := os.Getenv("API_SERVER")
	token := os.Getenv("TOKEN")
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if apiserver == "" {
		log.Fatal("API Server not defined")
	}

	if token == "" {
		log.Fatal("Token not defined")
	}

	config := Config{
		ApiServer:   apiserver,
		BearerToken: token,
		Insecure:    true}

	if debug {
		log.Printf("Accessing %v using %v\n", apiserver, token)
	}

	nodes, err := getPods(config)
	// nodes, err := getUnschedulable(config)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range nodes.Items {
		log.Printf("Node %v at %v with IP: %v", item.Metadata.Name, item.Status.HostIP, item.Status.PodIP)
	}

	// ports, error := getPorts(config)
	// if error != nil {
	// 	log.Fatal(error)
	// }

	// // extract node ip
	// for _, item := range nodes.Items {
	// 	for _, elem := range ports {
	// 		log.Printf("Node %v:%v", item.Metadata.Name, elem)
	// 	}
	// }

	// generate config for haproxy

	// update haproxy

}
