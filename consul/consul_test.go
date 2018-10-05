package consul

import (
	"fmt"
	"testing"

	"github.com/apex/log"
)

func TestWatchServiceTest(t *testing.T) {
	if url, err := WatchServiceURL("api-core"); err != nil {
		new(testing.T).Error(
			"For", "watch service test",
			"expected", "something",
			"got", "nothing",
		)
	} else {
		currentURL := *url
		log.Info(fmt.Sprintf("Initial URL value: %s", currentURL))
		for true {
			if currentURL != *url {
				currentURL = *url
				log.Info(fmt.Sprintf("url changed %s", currentURL))
			}
		}
	}

}
