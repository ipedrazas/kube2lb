package main

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

func writeToETCD(endpoints []string, key string, value string) (bool, error) {

	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return false, err
	}
	kapi := client.NewKeysAPI(c)

	resp, err := kapi.Set(context.Background(), key, value, nil)
	if err != nil {
		return false, err
	} else {
		log.Printf("Get is done. Metadata is %q\n", resp)
		return true, nil
	}

}
