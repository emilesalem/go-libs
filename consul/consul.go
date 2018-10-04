/*
 */
package consul

import (
	"fmt"

	"github.com/apex/log"
	consul "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/watch"
)

type ServiceInfo struct {
	Name string
	URL  string
}

func watchService(serviceInfo *ServiceInfo, client *consul.Client) error {
	plan, err := createServiceWatchPlan(serviceInfo)
	if err != nil {
		return err
	}
	plan.Handler = func(someInt uint64, result interface{}) {
		var results []*consul.ServiceEntry = result.([]*consul.ServiceEntry)
		selectNodeAddress(serviceInfo, results)
	}
	err = plan.RunWithClientAndLogger(client, nil)
	if err != nil {
		log.WithError(err).Error("error during service watch")
		return err
	}
	return nil
}

func createServiceWatchPlan(serviceInfo *ServiceInfo) (*consulWatch.Plan, error) {
	watchQuery := make(map[string]interface{})
	watchQuery["type"] = "service"
	watchQuery["service"] = serviceInfo.Name
	plan, err := consulWatch.Parse(watchQuery)
	if err != nil {
		log.WithError(err).Error("cant parse watch query parameters")
		return nil, err
	}
	return plan, nil
}

func selectNodeAddress(serviceInfo *ServiceInfo, nodes []*consul.ServiceEntry) {
	for _, n := range nodes {
		serviceAddress := n.Service.Address
		if len(serviceAddress) == 0 {
			serviceAddress = n.Node.Address
		}
		serviceAddress = fmt.Sprintf("%s:%v", serviceAddress, n.Service.Port)
		if serviceAddress != serviceInfo.URL {
			log.WithFields(log.Fields{
				"name": serviceInfo.Name,
				"url":  serviceAddress,
			}).Info("updating node")
			serviceInfo.URL = serviceAddress
			return
		}
	}
}

func WatchServiceURL(serviceName string) (*string, error) {
	serviceInfo := &ServiceInfo{serviceName, ""}
	client, err := consul.NewClient(&consul.Config{Address: ConsulAddress})
	if err != nil {
		log.WithError(err).Error("can't connect to consul")
		return nil, err
	}
	err = watchService(serviceInfo, client)
	if err != nil {
		return nil, err
	}
	return &serviceInfo.URL, nil
}
