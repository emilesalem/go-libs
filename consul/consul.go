// Package consul provides the means to get the changing value of a service URL
package consul

import (
	"errors"
	"fmt"
	stdLog "log"
	"math/rand"
	"os"
	"strconv"
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

var maxTime time.Duration
var consulClient *consulapi.Client

func init() {
	var err error
	consulClient, err = consulapi.NewClient(&consulapi.Config{Address: os.Getenv("CONSUL_HOST")})
	if err != nil {
		log.WithError(err).Error("cannot create consul client - panicking")
		panic(err)
	}
	v, _ := strconv.ParseInt(os.Getenv("INITIAL_VALUE_TIMEOUT_SECONDS"), 10, 0)
	maxTime = time.Duration(v)
	log.Info("consul client created")
}

// WatchService takes a service name and returns a pointer to a ServiceInfo holding the
// current URL of a healthy node (randomly chosen) of the service.
// The URL value will get updated as the service nodes change;
// the function will block until either of the following events occur:
// - the INITIAL_VALUE_TIMEOUT_SECONDS duration is elapsed
// - the service URL was resolved by consul
// if the INITIAL_VALUE_TIMEOUT_SECONDS duration is elapsed, an error is returned
func WatchService(serviceName string) (*ServiceInfo, error) {
	chCurrentValue := make(chan *ServiceInfo)
	chStop := make(chan bool)
	startWatch(serviceName, chCurrentValue, chStop)
	select {
	case serviceInfo := <-chCurrentValue:
		return serviceInfo, nil
	case <-time.After(maxTime * time.Second):
		msg := "consul watch timed out"
		log.Error(msg)
		chStop <- true
		close(chStop)
		return nil, errors.New(msg)
	}
}

func startWatch(serviceName string, c chan *ServiceInfo, chStop chan bool) {
	plan, err := createServiceWatchPlan(serviceName)
	if err != nil {
		log.WithError(err).Error("error creating service watch plan - panicking")
		panic(err)
	}
	serviceInfo := &ServiceInfo{serviceName, ""}
	var f bool
	plan.Handler = func(i uint64, result interface{}) {
		serviceInfo.URL = selectServiceEntryAddress(result.([]*consulapi.ServiceEntry))
		log.Info(fmt.Sprintf("updating %s service, new URL: %s", serviceInfo.Name, serviceInfo.URL))
		if !f {
			c <- serviceInfo
			close(c)
			f = !f
		}
	}
	go plan.RunWithClientAndLogger(consulClient, stdLog.New(os.Stderr, "", 1))
	go waitForStopSignal(plan, chStop)
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

func waitForStopSignal(plan *consulwatch.Plan, chStop chan bool) {
	<-chStop
	log.Info("stopping plan")
	plan.Stop()
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
