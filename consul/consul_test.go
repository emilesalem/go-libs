package consul

import (
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func TestConsul(t *testing.T) {
	serviceName := "test"
	servicePort := "12345"
	os.Setenv("SERVICE_NAME", serviceName)
	os.Setenv("SERVICE_PORT", servicePort)

	t.Run("Register service test", func(t *testing.T) {
		if err := RegisterLocalService(); err != nil {
			t.FailNow()
		}
		if nodes, _, err := consulClient.Catalog().Service(serviceName, "", nil); err != nil || len(nodes) == 0 {
			t.Fail()
		}
	})
	t.Run("Watch service test", func(t *testing.T) {
		serviceInfo, err := WatchService(serviceName)
		if err != nil {
			t.Error(
				"For", "watch service test",
				"expected", "service info",
				"got", "error",
			)
		}
		if serviceInfo == nil {
			t.FailNow()
		}

		previousInfo := *serviceInfo
		os.Setenv("SERVICE_PORT", "54321")
		if err := RegisterLocalService(); err != nil {
			t.FailNow()
		}
		if nodes, _, err := consulClient.Catalog().Service(serviceName, "", nil); err != nil || len(nodes) != 1 {
			t.Fail()
		}
		<-time.NewTimer(1 * time.Second).C
		if serviceInfo.URL == previousInfo.URL {
			t.Error(
				"For", "updated service",
				"expected", "127.0.0.1:54321",
				"got", serviceInfo.URL,
			)
		}
	})
}
