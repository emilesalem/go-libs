package health

import (
	"net/http"
	"testing"
)

func TestHealthyHealthChecks(t *testing.T) {
	resetHealthcheckTest()
	for _, healthChecker := range createHealthyHealthcheckers() {
		AddHealthCheck(healthChecker)
	}
	Handler(&healthyHealthcheckTestWriter{t}, nil)
	allHealthchecksCalled(t)
}

type healthyHealthcheckTestWriter struct {
	test *testing.T
}

func (t *healthyHealthcheckTestWriter) WriteHeader(status int) {
	if status != http.StatusOK {
		new(testing.T).Error(
			"For", "healthy healthcheck",
			"expected", http.StatusOK,
			"got", status,
		)
	}
}

func (t *healthyHealthcheckTestWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (t *healthyHealthcheckTestWriter) Header() http.Header {
	return make(http.Header)
}
