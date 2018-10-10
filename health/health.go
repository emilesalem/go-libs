package health

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Workiva/go-datastructures/futures"
	"github.com/apex/log"
	"github.com/emilesalem/go-libs/env"
)

var healthChecks []func(*chan interface{}, *futures.Future)
var healthTimeout time.Duration

func init() {
	if t, err := strconv.Atoi(env.Get("HEALTH_TIMEOUT")); err != nil {
		log.WithError(err).Fatal("could not parse HEALTH_TIMEOUT value to int")
		panic(err)
	} else {
		healthTimeout = time.Duration(t) * time.Second
	}
}

func AddHealthCheck(f func() error) {
	healthChecks = append(healthChecks, func(c *chan interface{}, future *futures.Future) {
		if !future.HasResult() {
			*c <- f()
		}
	})
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	errorC := make(chan interface{})
	startHealthchecks(&errorC)
	w.Header().Set("Content-Type", "application/json")
	if prognosisNegative(&errorC) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{ "status" : "OUT_OF_SERVICE" }`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "status" : "UP" }`))
	}
}

func startHealthchecks(errorC *chan interface{}) {
	for _, f := range healthChecks {
		go f(errorC, futures.New(*errorC, healthTimeout))
	}
}

func prognosisNegative(errorC *chan interface{}) bool {
	healthCount := 0
	for err := range *errorC {
		if err != nil {
			log.WithError(err.(error)).Warn("failed health check")
			return true
		}
		healthCount++
		if healthCount == len(healthChecks) {
			close(*errorC)
		}
	}
	return false
}
