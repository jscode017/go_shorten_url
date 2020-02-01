package main

import (
	"log"
	"time"
)

func ConvetStrToTime(timeStr string) time.Time {
	layout := "2006-01-02 15:04:05 -0700 MST"
	convertedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Fatal(err)
	}

	return convertedTime

}
func GetElapsedTime(timeStr string) int64 { //return remain seconds
	timeNow := time.Now()
	CreationTime := ConvetStrToTime(timeStr)

	elapsedTimeInSecs := int64(timeNow.Sub(CreationTime).Seconds())
	return elapsedTimeInSecs
}

func ConvertTimeToStr(t time.Time) string {
	return t.UTC().String()
}
