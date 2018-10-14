package discovery

import (
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {

	serviceName := "test"
	servicePort := 12345

	config := Config{
		"127.0.0.1:8500",
		true,
		serviceName,
		servicePort,
	}

	w := MakeDiscoveryService(config)
	d := w.(discoveryService)
	t.Run("Register service test", func(t *testing.T) {
		if err := d.registerLocalService(serviceName, servicePort); err != nil {
			t.FailNow()
		}
		if nodes, _, err := d.consulClient.Catalog().Service(serviceName, "", nil); err != nil || len(nodes) == 0 {
			t.Fail()
		}
	})
	t.Run("Watch service test", func(t *testing.T) {
		serviceInfoC := d.watchService(serviceName)
		if serviceInfoC == nil {
			t.FailNow()
		}

		previousInfo := <-serviceInfoC
		if err := d.registerLocalService(serviceName, 54321); err != nil {
			t.FailNow()
		}
		if nodes, _, err := d.consulClient.Catalog().Service(serviceName, "", nil); err != nil || len(nodes) != 1 {
			t.Fail()
		}
		<-time.NewTimer(1 * time.Second).C
		currentInfo := <-serviceInfoC
		if currentInfo.URL == previousInfo.URL {
			t.Error(
				"For", "updated service",
				"expected", "127.0.0.1:54321",
				"got", currentInfo.URL,
			)
		}
	})
}
