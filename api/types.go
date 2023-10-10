package api

import (
	"encoding/json"
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

var PodsMap = make(map[string]Pod)
var PodMetricsMap = make(map[string]PodMetrics)

type PodMetrics struct {
	PodName string
	Events  []Event
}

type Message struct {
	Source  string          `json:"source"`
	Target  string          `json:"target"`
	Content string          `json:"content"`
	Context json.RawMessage `json:"context"` // RawMessage here! (not a string)
}
