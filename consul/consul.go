/*
 */
package consul

import (
	"fmt"
	stdLog "log"
	"math/rand"
	"os"
	"time"

	"github.com/apex/log"
	consul "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/watch"
)

type ServiceInfo struct {
	Name string
	URL  string
}

func watchService(serviceInfo *ServiceInfo, client *consul.Client) {
	plan, err := createServiceWatchPlan(serviceInfo)
	if err != nil {
		log.WithError(err).Error("error creating service watch")
		return
	}
	plan.Handler = func(someInt uint64, result interface{}) {
		results := result.([]*consul.ServiceEntry)
		selectServiceEntryNodeAddress(results)
		log.Info(fmt.Sprintf("updating service address: %s", serviceInfo.URL))
	}
	err = plan.RunWithClientAndLogger(client, stdLog.New(os.Stderr, "", 1))
	if err != nil {
		log.WithError(err).Error("error starting service watch")
	}
}

func fetchInitialServiceURL(serviceName string, catalog *consul.Catalog) (string, error) {
	nodes, _, err := catalog.Service(serviceName, "", nil)
	if err != nil {
		log.WithError(err).Error("can't get services from catalog")
		return "", err
	}
	if len(nodes) <= 0 {
		log.WithField("name", serviceName).Error("service not found")
		return "", err
	}
	return selectCatalogEntryNodeAddress(nodes), nil
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

func selectServiceEntryNodeAddress(nodes []*consul.ServiceEntry) string {
	var result string = ""
	if len(nodes) > 0 {
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := nodes[r1.Intn(len(nodes))]
		serviceAddress := n.Service.Address
		if len(serviceAddress) == 0 {
			serviceAddress = n.Node.Address
		}
		result = fmt.Sprintf("%s:%v", serviceAddress, n.Service.Port)
	}
	return result
}

func selectCatalogEntryNodeAddress(nodes []*consul.CatalogService) string {
	var result string = ""
	if len(nodes) > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := nodes[r.Intn(len(nodes))]
		serviceAddress := n.ServiceAddress
		if len(serviceAddress) == 0 {
			serviceAddress = n.Address
		}
		result = fmt.Sprintf("%s:%v", serviceAddress, n.ServicePort)
	}
	return result
}

func WatchServiceURL(serviceName string) (*ServiceInfo, error) {
	serviceInfo := &ServiceInfo{serviceName, ""}
	client, err := consul.NewClient(&consul.Config{Address: ConsulAddress})
	if err != nil {
		log.WithError(err).Error("can't connect to consul")
		return nil, err
	}
	if serviceInfo.URL, err = fetchInitialServiceURL(serviceInfo.Name, client.Catalog()); err != nil {
		return nil, err
	}
	go watchService(serviceInfo, client)
	return serviceInfo, nil
}
