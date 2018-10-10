package health

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Workiva/go-datastructures/futures"
)

func resetHealthcheckTest() {
	healthChecks = []func(*chan interface{}, *futures.Future){}
	healthCheckCalls = [4]bool{}
}

var healthCheckCalls [4]bool

func createHealthyHealthcheckers() []func() error {
	var healthCheckers []func() error
	for i := 0; i < 4; i++ {
		index := i
		healthCheckers = append(healthCheckers, func() error {
			duration := time.Duration(rand.New(rand.NewSource(0)).Intn(300))
			<-time.NewTimer(duration * time.Millisecond).C
			healthCheckCalls[index] = true
			return nil
		})
	}
	return healthCheckers
}

func allHealthchecksCalled(t *testing.T) {
	for _, value := range healthCheckCalls {
		if !value {
			t.Error(
				"For", "healthcheck called",
				"expected", true,
				"got", value,
			)
		}
	}
}
