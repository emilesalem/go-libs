package consul

import (
	"fmt"
	"testing"

	"github.com/apex/log"
)

func TestWatchServiceTest(t *testing.T) {
	urlUGC, err := WatchServiceURL("api-ugc")
	if err != nil {
		t.Error(
			"For", "watch service test",
			"expected", "service info",
			"got", "error",
		)
		t.FailNow()
	}
	urlSSO, err := WatchServiceURL("api-authentication-sso")
	if err != nil {
		t.Error(
			"For", "watch service test",
			"expected", "service info",
			"got", "error",
		)
		t.FailNow()
	}

	currentURLUGC := *urlUGC
	log.Info(fmt.Sprintf("Initial ugc URL value: %s", currentURLUGC.URL))
	currentURLSSO := *urlSSO
	log.Info(fmt.Sprintf("Initial sso URL value: %s", currentURLSSO.URL))
	for true {
		if currentURLUGC.URL != urlUGC.URL {
			currentURLUGC = *urlUGC
			log.Info(fmt.Sprintf("ugc url changed %s", currentURLUGC.URL))
		}
		if currentURLSSO.URL != urlSSO.URL {
			currentURLSSO = *urlSSO
			log.Info(fmt.Sprintf("ugc sso changed %s", currentURLSSO.URL))
		}
	}

}
