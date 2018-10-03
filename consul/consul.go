package consul

import (
	"fmt"
	"time"

	"github.com/apex/log"
	consul "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/watch"
)

const ServiceDiscoveryRefreshDuration = 10 * time.Second

type ServiceInfo struct {
	Name string
	URL  string
}

/*
see:
https://www.consul.io/docs/agent/watches.html
https://github.com/hashicorp/consul/blob/master/watch/plan.go
https://github.com/hashicorp/consul/blob/master/watch/watch.go#L116
https://gowalker.org/github.com/hashicorp/consul/watch#WatcherFunc
https://godoc.org/github.com/hashicorp/consul/watch#Plan
https://godoc.org/github.com/hashicorp/consul/api#CatalogService

*/
func watchService(serviceName string, client *consul.Client) {
	watchQuery := make(map[string]interface{})
	watchQuery["type"] = "service"
	watchQuery["service"] = serviceName
	plan, err := consulWatch.Parse(watchQuery)
	if err != nil {
		log.WithError(err).Error("cant parse watch query parameters")
	} else {
		log.Info(fmt.Sprintf("plan info! %v", plan))
	}
	plan.Handler = func (someInt uint64, result consul.CatalogService) {
		log.Info(fmt.Sprintf("%v", result.Address))
	}
	plan.RunWithClientAndLogger(client, nil)
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
		if url := fmt.Sprintf("http://%v:%v", s.ServiceAddress, s.ServicePort); url != serviceInfo.URL {
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
		watchService(serviceInfo.Name, client)
	}
}

func WatchService(serviceName string) *string {
	serviceInfo := &ServiceInfo{serviceName, ""}
	discover(serviceInfo) // invoke first tick
	// ticker := time.NewTicker(ServiceDiscoveryRefreshDuration)
	// go func() {
	// 	for range ticker.C {
	// 		discover(serviceInfo)
	// 	}
	// }()
	return &serviceInfo.URL
}
