package metrics

import (
	v1 "k8s.io/api/core/v1"
	"strings"
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
type PodMetrics struct {
	PodName string
	Events  []Event
}

var PodMetricsMap = make(map[string]PodMetrics)

// fmt.Printf("%s Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
//
//	time.Now(), event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)

type PodMetricsInterface interface {
	GetSchedulingTime() (int, error)
	GetPrepTime() (int, error)
	GetRunningTime() (int, error)
}

func NewPodMetrics(podK8s *v1.Pod, typeEvent string, phase string) (*PodMetrics, error) {
	podT := Pod{
		PodName:  podK8s.ObjectMeta.Name,
		NodeName: podK8s.Spec.NodeName,
		NodeIP:   podK8s.Status.HostIP,
		PodIp:    podK8s.Status.PodIP,
	}
	metricEvent := Event{
		Type:      typeEvent,
		Pod:       podT,
		Phase:     phase,
		Timestamp: time.Now(),
	}

	podMetric, exists := PodMetricsMap[podK8s.ObjectMeta.Name]
	if !exists {
		podMetric = PodMetrics{
			PodName: podK8s.ObjectMeta.Name,
		}
	}
	podMetric.Events = append(podMetric.Events, metricEvent)
	PodMetricsMap[podK8s.ObjectMeta.Name] = podMetric
	return &podMetric, nil
}

func (p *PodMetrics) GetSchedulingTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var schedulingEvent []Event
	for _, event := range p.Events {
		if !strings.EqualFold(event.Phase, "Pending") {
			continue
		}

		if strings.EqualFold(event.Type, "ADDED") {
			schedulingEvent = append(schedulingEvent, event)
		} else if strings.EqualFold(event.Type, "MODIFIED") && len(event.Pod.NodeName) > 0 &&
			len(event.Pod.PodIp) == 0 && len(event.Pod.NodeIP) == 0 {
			schedulingEvent = append(schedulingEvent, event)
			break
		}
	}

	if len(schedulingEvent) < 2 {
		return 0
	}

	duration := schedulingEvent[1].Timestamp.Sub(schedulingEvent[0].Timestamp)

	return int(duration.Milliseconds())
}

func (p *PodMetrics) GetPrepTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var events []Event
	for _, event := range p.Events {
		if !strings.EqualFold(event.Type, "MODIFIED") ||
			len(event.Pod.NodeName) == 0 {
			continue
		}

		if len(event.Pod.PodIp) == 0 && len(event.Pod.NodeIP) == 0 {
			events = append(events, event)
		} else if len(event.Pod.PodIp) > 0 && len(event.Pod.NodeIP) > 0 {
			events = append(events, event)
			break
		}
	}

	if len(events) < 2 {
		return 0
	}

	duration := events[1].Timestamp.Sub(events[0].Timestamp)
	return int(duration.Milliseconds())
}

func (p *PodMetrics) GetRunningTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var events []Event
	for _, event := range p.Events {
		if strings.EqualFold(event.Type, "DELETED") || len(event.Pod.NodeName) == 0 {
			continue
		}

		if len(event.Pod.PodIp) > 0 && len(event.Pod.NodeIP) > 0 {
			events = append(events, event)
			if len(events) == 2 {
				break
			}
		}
	}

	if len(events) < 2 {
		return 0
	}

	duration := events[1].Timestamp.Sub(events[0].Timestamp)
	return int(duration.Milliseconds())
}
