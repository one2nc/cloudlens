package gcp

import (
	"fmt"
	"strings"
	"time"

	"github.com/one2nc/cloudlens/internal/config"
)

func GetLocalTime(timestamp string) (string, error) {

	launchTime, err := time.Parse("2006-01-02T15:04:05.999-07:00", timestamp)
	if err != nil {
		return "", fmt.Errorf("Error parsing timestamp : ", err)
	}
	localZone, err := config.GetLocalTimeZone()
	if err != nil {
		return "", fmt.Errorf("Error loading local timezone: ", err)
	}
	loc, _ := time.LoadLocation(localZone)
	IST := launchTime.In(loc)
	return IST.Format("Mon Jan _2 15:04:05 2006"), nil
}

func GetResourceFromURL(url string) string {
	splittedURL := strings.Split(url, "/")
	resource := splittedURL[len(splittedURL)-1]
	return resource
}
