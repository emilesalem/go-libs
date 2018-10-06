// Package consul provides the means to get the changing value of a service URL
package consul

import (
	"fmt"
	stdLog "log"
	"math/rand"
	"os"
	"time"

	"github.com/apex/log"
	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/watch"
)

//ServiceInfo holds the name and the URL of a service
type ServiceInfo struct {
	Name string
	URL  string
}

//WatchServiceURL returns a pointer to a ServiceInfo holding the current URL of the service
//this URL value will get updated as the service nodes change
func WatchServiceURL(serviceName string) (*ServiceInfo, error) {
	serviceInfo := &ServiceInfo{serviceName, ""}
	client, err := consulapi.NewClient(&consulapi.Config{Address: ConsulAddress})
	if err != nil {
		log.WithError(err).Error("can't connect to consul")
		return nil, err
	}
	if serviceInfo.URL, err = fetchCurrentServiceURL(serviceInfo.Name, *client); err != nil {
		return nil, err
	}
	go watchService(serviceInfo, *client)
	return serviceInfo, nil
}

func fetchCurrentServiceURL(serviceName string, client consulapi.Client) (string, error) {
	nodes, _, err := client.Health().Service(serviceName, "", true, &consulapi.QueryOptions{AllowStale: false})
	if err != nil {
		log.WithError(err).Error("can't get services from consul")
		return "", err
	}
	if len(nodes) <= 0 {
		log.WithField("name", serviceName).Error(fmt.Sprintf("%s service not found", serviceName))
		return "", err
	}
	return selectServiceEntryAddress(nodes), nil
}

func selectServiceEntryAddress(nodes []*consulapi.ServiceEntry) string {
	var result string
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

func watchService(serviceInfo *ServiceInfo, client consulapi.Client) {
	plan, err := createServiceWatchPlan(serviceInfo.Name)
	if err != nil {
		log.WithError(err).Error("error creating service watch plan")
		return
	}
	plan.Handler = func(someInt uint64, result interface{}) {
		serviceInfo.URL = selectServiceEntryAddress(result.([]*consulApi.ServiceEntry))
		log.Info(fmt.Sprintf("updating %s service, new URL: %s", serviceInfo.Name, serviceInfo.URL))
	}
	err = plan.RunWithClientAndLogger(&client, stdLog.New(os.Stderr, "", 1))
	if err != nil {
		log.WithError(err).Error("error starting service watch")
	}
}

func createServiceWatchPlan(serviceName string) (*consulwatch.Plan, error) {
	watchQuery := make(map[string]interface{})
	watchQuery["type"] = "service"
	watchQuery["service"] = serviceName
	watchQuery["passingonly"] = true
	plan, err := consulwatch.Parse(watchQuery)
	if err != nil {
		log.WithError(err).Error("cant parse watch query parameters")
		return nil, err
	}
	return plan, nil
}
