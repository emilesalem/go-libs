// Package consul provides the means to get the changing value of a service URL.
// If we detect 'development' environment we register the service
// using SERVICE_NAME and SERVICE_PORT envars and localhost address.
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
	"github.com/emilesalem/go-libs/env"
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
	initConsulClient()
	initMaxTime()
	if env.GetOpt("ENVIRONMENT") == "dev" {
		RegisterLocalService()
	}
}

func initConsulClient() {
	if c, err := consulapi.NewClient(&consulapi.Config{Address: env.Get("CONSUL_HOST")}); err != nil {
		log.WithError(err).Fatal("cannot create consul client")
		panic(err)
	} else {
		consulClient = c
		log.Info("consul client created")
	}
}

func initMaxTime() {
	if v, err := strconv.Atoi(env.Get("INITIAL_WATCH_TIMEOUT_SECONDS")); err != nil {
		log.WithError(err).Fatal("INITIAL_WATCH_TIMEOUT_SECONDS cannot be parsed to int")
		panic(err)
	} else {
		maxTime = time.Duration(v)
	}
}

// WatchService accepts a service name and returns a ServiceInfo pointer holding the
// current URL of a random healthy service node.
// The URL value will get updated as the service nodes change;
// the function will block until either of the following events occur:
// - the maxTime duration is elapsed
// - the service URL was resolved by consul
// if the maxTime duration is elapsed an error is returned
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
		log.WithError(err).Fatal("error creating service watch plan")
		panic(err)
	}
	serviceInfo := &ServiceInfo{serviceName, ""}
	var f bool
	plan.Handler = func(i uint64, result interface{}) {
		if selectedNode := selectServiceEntryAddress(result.([]*consulapi.ServiceEntry)); len(selectedNode) > 0 {
			serviceInfo.URL = selectedNode
			log.Infof("updating %s service, new URL: %s", serviceInfo.Name, serviceInfo.URL)
			if !f {
				c <- serviceInfo
				close(c)
				f = !f
			}
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
	log.Info("stopping watch plan")
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
