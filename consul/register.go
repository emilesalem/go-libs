package consul

import (
	"github.com/apex/log"
	consulapi "github.com/hashicorp/consul/api"
)

//RegisterService registers a service to Consul
func RegisterService(name string, address string, port int) error {
	log.Infof("registrating service %s at %s:$v", name, address, port)
	if err := consulClient.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		Name:    name,
		Address: address,
		Port:    port,
	}); err != nil {
		log.WithError(err).Errorf("failed to register service %s", name)
		return err
	}
	return nil
}
