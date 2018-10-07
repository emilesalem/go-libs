package consul

import (
	"fmt"
	"testing"
	"time"

	"github.com/apex/log"
	_ "github.com/joho/godotenv/autoload"
)

func TestWatchServiceTest(t *testing.T) {
	ugcInfo, err := WatchService("api-ugc")
	if err != nil {
		t.Error(
			"For", "watch service test",
			"expected", "service info",
			"got", "error",
		)
	}
	if ugcInfo == nil {
		t.FailNow()
	}

	currentURLUGC := *ugcInfo
	log.Info(fmt.Sprintf("Initial ugc URL value: %s", currentURLUGC.URL))
	for range time.NewTicker(2 * time.Second).C {
		if currentURLUGC.URL != ugcInfo.URL {
			currentURLUGC = *ugcInfo
			log.Info(fmt.Sprintf("ugc url changed %s", currentURLUGC.URL))
		}
	}
}
