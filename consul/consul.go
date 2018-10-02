package consul

import (
	consul "github.com/hashicorp/consul/api"
	"github.com/apex/log"
	"strconv"
	"time"
)

const ServiceDiscoveryRefreshDuration = 10 * time.Second

type ServiceInfo struct { 
	Name string
	URL  string
}

func findService(catalog *consul.Catalog, serviceInfo *ServiceInfo) {
	services, _, err := catalog.Service(serviceInfo.Name, "", nil)

	if err != nil {
		log.WithError(err).Error("can't get services from catalog")
		return 
	}

	if len(services) <= 0 {
		log.WithField("name", serviceInfo.Name).Error("service not found")
		return
	}

	for _, s := range services {
		if url := "http://" + s.ServiceAddress + ":" + strconv.Itoa(s.ServicePort); url != serviceInfo.URL {
			log.WithFields(log.Fields{
				"name": serviceInfo.Name,
				"url":  url,
			}).Info("updating node")
			serviceInfo.URL = url
			return
		}
	}
}

func discover(serviceInfo *ServiceInfo) {
	client, err := consul.NewClient(&consul.Config{Address: ConsulAddress})
	if err != nil {
		log.WithError(err).Error("can't connect to consul")
	} else {
		findService(client.Catalog(), serviceInfo)
	}
}

func WatchService(serviceName string) *string{
	serviceInfo := &ServiceInfo{serviceName, ""}
	discover(serviceInfo) // invoke first tick
	ticker := time.NewTicker(ServiceDiscoveryRefreshDuration)
	go func() {
		for range ticker.C {
			discover(serviceInfo)
		}
	}()
	return &serviceInfo.URL
}