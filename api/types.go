package api

import (
	"log"
	"time"
)

type Pod struct {
	PodName     string
	NodeName    string
	NodeIP      string
	PodIp       string
	Annotations map[string]string
}
type Event struct {
	Type      string
	Pod       Pod
	Phase     string
	Timestamp time.Time
}

var Debug, Trace bool
var InfoLog, DebugLog, ErrorLog, TraceLog *log.Logger
var PodsMap = new(Map[string, Pod])
var PodMetricsMap = make(map[string]PodMetrics)

type PodMetrics struct {
	PodName string
	Events  []Event
}
