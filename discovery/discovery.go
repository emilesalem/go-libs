// Package discovery provides the means to get the changing value of a service URL.
package discovery

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

//ServiceInfo values hold the name and the URL of a service
type ServiceInfo struct {
	Name string
	URL  string
}

//ServiceWatcher is destined to be used to discover other services using their name
type ServiceWatcher interface {
	watchService(string) <-chan ServiceInfo
}

type discoveryService struct {
	consulClient *consulapi.Client
}

// MakeDiscoveryService creates a new discovery service
func MakeDiscoveryService(config Config) ServiceWatcher {
	consulClient, err := consulapi.NewClient(&consulapi.Config{Address: config.consulAddress})
	if err != nil {
		log.WithError(err).Fatal("cannot create consul client")
		panic(err)
	}
	log.Info("consul client created")
	discoveryService := discoveryService{consulClient}
	if config.localRegistration {
		discoveryService.registerLocalService(config.serviceName, config.servicePort)
	}
	return discoveryService
}

// WatchService accepts a service name and returns a ServiceInfo receiving channel;
// The ServiceInfo sent through the channel will hold the URL of a random healthy service node.
func (s discoveryService) watchService(name string) <-chan ServiceInfo {
	c := make(chan ServiceInfo)
	s.startWatch(name, c)
	return c
}

func (s discoveryService) startWatch(serviceName string, c chan ServiceInfo) {
	plan, err := createServiceWatchPlan(serviceName)
	if err != nil {
		log.WithError(err).Fatal("error creating service watch plan")
		panic(err)
	}
	serviceInfo := ServiceInfo{serviceName, ""}
	plan.Handler = func(i uint64, result interface{}) {
		if selectedNode := selectServiceEntryAddress(result.([]*consulapi.ServiceEntry)); len(selectedNode) > 0 {
			serviceInfo.URL = selectedNode
			log.Infof("updating %s service, new URL: %s", serviceInfo.Name, serviceInfo.URL)
			c <- serviceInfo
		}
	}
	go plan.RunWithClientAndLogger(s.consulClient, stdLog.New(os.Stderr, "", 1))
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
