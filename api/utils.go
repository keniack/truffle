package api

import (
	"fmt"
	"os"
	"time"
)

type LogWriter struct {
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.9999999Z") + " [INFO] " + string(bytes))
}

func ExistsInMap(myMap *Map[string, Pod], key string) bool {
	_, exists := myMap.Load(key)
	return exists
}

func ExistsAndHasPodIp(podMap *Map[string, Pod], key string) bool {
	pod, exists := podMap.Load(key)
	if exists && len(pod.PodIp) > 0 {
		return true
	}
	return false
}

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
}
