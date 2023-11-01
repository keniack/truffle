package metrics

import (
	v1 "k8s.io/api/core/v1"
	"polaris/truffle/pkg/common"
	"strings"
	"time"
)

// fmt.Printf("%s Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
//
//	time.Now(), event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)

/*type PodMetricsInterface interface {
	GetSchedulingTime() (int, error)
	GetPrepTime() (int, error)
	GetRunningTime() (int, error)
}

*/

func NewPodMetrics(podK8s *v1.Pod, typeEvent string, phase string) (*common.PodMetrics, error) {
	podT := common.Pod{
		PodName:  podK8s.ObjectMeta.Name,
		NodeName: podK8s.Spec.NodeName,
		NodeIP:   podK8s.Status.HostIP,
		PodIp:    podK8s.Status.PodIP,
	}
	metricEvent := common.Event{
		Type:      typeEvent,
		Pod:       podT,
		Phase:     phase,
		Timestamp: time.Now(),
	}

	podMetric, exists := common.PodMetricsMap[podK8s.ObjectMeta.Name]
	if !exists {
		podMetric = common.PodMetrics{
			PodName: podK8s.ObjectMeta.Name,
		}
	}
	podMetric.Events = append(podMetric.Events, metricEvent)
	common.PodMetricsMap[podK8s.ObjectMeta.Name] = podMetric
	return &podMetric, nil
}

func (p *pkg.PodMetrics) Exists() bool {
	_, exists := common.PodMetricsMap[p.PodName]
	return exists
}

func (p *pkg.PodMetrics) GetSchedulingTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var schedulingEvent []common.Event
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

func (p *pkg.PodMetrics) GetPrepTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var events []common.Event
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

func (p *pkg.PodMetrics) GetRunningTime() int {

	if len(p.Events) == 0 {
		return 0
	}
	var events []common.Event
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
