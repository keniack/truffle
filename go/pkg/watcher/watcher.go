package watcher

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
	"polaris/truffle/pkg/common"
	"strings"
	"sync"
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

func GetPodIpForName(target string, podIpChannel chan<- string, wg *sync.WaitGroup) {
	if common.Debug {
		common.DebugLog.Printf("starting watcher")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer wg.Done()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}

	// watch future changes to Pods
	watcher, err := WatcherClient().CoreV1().Pods("default").Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()
out:
	for event := range ch {
		pod, _ := event.Object.(*v1.Pod)
		if common.Trace {
			common.TraceLog.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
				event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)
		}
		switch {
		case (event.Type == watch.Modified || event.Type == watch.Added) && strings.HasPrefix(pod.ObjectMeta.Name, target) && len(pod.Status.PodIP) > 0:
			{
				if common.Debug {
					common.DebugLog.Println("Adding podIp to channel")
				}
				podIpChannel <- pod.Status.PodIP
				break out

			}
		}
	}
	if common.Debug {
		common.DebugLog.Println("Finished watcher")
	}

}

func GetNodeIpForName(target string, nodeIpChannel chan<- string, wg *sync.WaitGroup) {
	if common.Debug {
		common.DebugLog.Printf("starting watcher")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer wg.Done()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}

	// watch future changes to Pods
	watcher, err := WatcherClient().CoreV1().Pods("default").Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()

out:
	for event := range ch {
		pod, _ := event.Object.(*v1.Pod)
		if common.Trace {
			common.TraceLog.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
				event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)
		}

		switch {
		case (event.Type == watch.Modified || event.Type == watch.Added) && strings.HasPrefix(pod.ObjectMeta.Name, target) && len(pod.Status.HostIP) > 0:
			{
				if common.Debug {
					common.DebugLog.Println("Adding nodeIp to channel")
				}
				nodeIpChannel <- pod.Status.HostIP
				break out
			}
		}
	}
	if common.Debug {
		common.DebugLog.Println("Finished watcher")
	}

}
