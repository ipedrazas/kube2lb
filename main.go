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
	base := os.Getenv("BASE")
	etcd := os.Getenv("ETCD_ENDPOINTS")

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

	// Case NodePorts
	// get nodes
	nodes, err := getNodes(config)

	if err != nil {
		log.Fatal(err)
	}
	// get nodeports
	ports, err := getPorts(config)
	if err != nil {
		log.Fatal(err)
	}

	endpoints := []string{etcd}
	for _, node := range nodes {
		for _, port := range ports {
			key := fmt.Sprintf("%v/%v", base, node)
			writeToETCD(endpoints, key, strconv.Itoa(port))
		}
	}

	// Case Direct Access to Pods
	// get endpoints by service

}
