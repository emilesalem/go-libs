package discovery

import (
	"github.com/apex/log"
	consulapi "github.com/hashicorp/consul/api"
)

//registerLocalService registers a local service to Consul
func (s discoveryService) registerLocalService(name string, port int) error {
	log.Infof("registering local service %s on port %v", name, port)
	if err := s.consulClient.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		Name:    name,
		Address: "127.0.0.1",
		Port:    port,
	}); err != nil {
		log.WithError(err).Errorf("failed to register service %s", name)
		return err
	}
	return nil
}
