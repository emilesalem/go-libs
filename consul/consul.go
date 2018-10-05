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

func watchService(serviceInfo *ServiceInfo, client *consul.Client, c chan *ServiceInfo) {
	plan, err := createServiceWatchPlan(serviceInfo)
	if err != nil {
		c <- nil
	}
	initialValue := true
	plan.Handler = func(someInt uint64, result interface{}) {
		results := result.([]*consul.ServiceEntry)
		selectNodeAddress(serviceInfo, results)
		log.Info(fmt.Sprintf("updating service address: %s", serviceInfo.URL))
		if initialValue {
			c <- serviceInfo
			close(c)
			initialValue = false
		}
	}
	err = plan.RunWithClientAndLogger(client, stdLog.New(os.Stderr, "", 1))
	if err != nil {
		log.WithError(err).Error("error during service watch")
		c <- nil
	}
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
	if len(nodes) > 0 {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		n := nodes[r1.Intn(len(nodes))]
		serviceAddress := n.Service.Address
		if len(serviceAddress) == 0 {
			serviceAddress = n.Node.Address
		}
		serviceInfo.URL = fmt.Sprintf("%s:%v", serviceAddress, n.Service.Port)
	}
}

func WatchServiceURL(serviceName string) (*string, error) {
	serviceInfo := &ServiceInfo{serviceName, ""}
	client, err := consul.NewClient(&consul.Config{Address: ConsulAddress})
	if err != nil {
		log.WithError(err).Error("can't connect to consul")
		return nil, err
	}
	initialValueChan := make(chan *ServiceInfo)
	go watchService(serviceInfo, client, initialValueChan)
	serviceInfo = <-initialValueChan
	return &serviceInfo.URL, nil
}
