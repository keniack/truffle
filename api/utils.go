package api

import (
	"fmt"
	"time"
)

type LogWriter struct {
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.9999999Z") + " [INFO] " + string(bytes))
}

func ExistsInMap[V any](myMap map[string]V, key string) bool {
	_, exists := myMap[key]
	return exists
}

func ExistsAndHasPodIp(podMap map[string]Pod, key string) bool {
	pod, exists := podMap[key]
	if exists && len(pod.PodIp) > 0 {
		return true
	}
	return false
}
