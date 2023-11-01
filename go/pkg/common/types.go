package common

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

var RedisIP string
var AwsAccessKey string
var AwsSecretKey string
var ComMode string
var IncomingPodPort string
var Debug, Trace bool
var InfoLog, DebugLog, ErrorLog, TraceLog *log.Logger

// var PodsMap = new(utils.Map[string, Pod])
var PodMetricsMap = make(map[string]PodMetrics)

type PodMetrics struct {
	PodName string
	Events  []Event
}
