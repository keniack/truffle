package api

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WatcherClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func GetPodIpForName(target string) string {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}

	// watch future changes to Pods
	watcher, err := WatcherClient().CoreV1().Pods("default").Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()
	for event := range ch {
		pod, _ := event.Object.(*v1.Pod)
		if Trace {
			TraceLog.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
				event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)
		}
		switch {
		case (event.Type == watch.Modified || event.Type == watch.Added) && strings.HasPrefix(pod.ObjectMeta.Name, target) && len(pod.Status.PodIP) > 0:
			{
				return pod.Status.PodIP

			}
		}

	}
	return ""
}

func GetNodeIpForName(target string) string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}

	// watch future changes to Pods
	watcher, err := WatcherClient().CoreV1().Pods("default").Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()

	for event := range ch {
		pod, _ := event.Object.(*v1.Pod)
		if Trace {
			TraceLog.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
				event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)
		}
		switch {
		case (event.Type == watch.Modified || event.Type == watch.Added) && strings.HasPrefix(pod.ObjectMeta.Name, target) && len(pod.Status.HostIP) > 0:
			{
				return pod.Status.HostIP
			}
		}

	}
	return ""
}
