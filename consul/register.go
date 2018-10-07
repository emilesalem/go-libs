package consul

import (
	"strconv"

	"github.com/apex/log"
	"github.com/emilesalem/go-libs/env"
	consulapi "github.com/hashicorp/consul/api"
)

//RegisterLocalService registers a local service to Consul
func RegisterLocalService() error {
	port, _ := strconv.Atoi(env.Get("SERVICE_PORT"))
	name := env.Get("SERVICE_NAME")
	log.Infof("registering local service %s on port %v", name, port)
	if err := consulClient.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		Name:    name,
		Address: "127.0.0.1",
		Port:    port,
	}); err != nil {
		log.WithError(err).Errorf("failed to register service %s", name)
		return err
	}
	return nil
}
