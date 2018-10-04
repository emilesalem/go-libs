package consul

import (
	"fmt"
	"testing"

	"github.com/apex/log"
)

func TestWatchServiceTest(t *testing.T) {
	if url, err := WatchServiceURL("api-core"); err != nil || url == nil {
		new(testing.T).Error(
			"For", "watch service test",
			"expected", "something",
			"got", "nothing",
		)
	} else {
		log.Info(fmt.Sprintf("test success, got URL: %s", *url))
	}

}
