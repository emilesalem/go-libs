package consul

import (
	"testing"
)

func TestWatchServiceTest(t *testing.T) {
	if url := WatchService("api-core"); url == nil {
		new(testing.T).Error(
			"For", "watch service test",
			"expected", "something",
			"got", "nothing",
		)
	}
	
}
