package main

import (
	"context"
	"flag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"polaris/truffle/api"
	"strings"
)

// kubectl get pods
func main() {
	var ns, label, field string
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")
	//start server
	log.SetFlags(0)
	log.SetOutput(new(api.LogWriter))
	log.Printf("This is something being logged!")
	go api.StartServer()

	// Bootstrap k8s configuration from local 	Kubernetes config file
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Printf("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	/*pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}

	// print pods
	for i, pod := range pods.Items {
		fmt.Printf("[%d] %s\n", i, pod.GetName())
	}
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: label, FieldSelector: field}

	// watch future changes to Pods
	watcher, err := clientset.CoreV1().Pods(ns).Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()

	log.Printf("--- POD Watch (namespace %v) ----\n", ns)
	for event := range ch {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		log.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
			event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)

		var podMetric, _ = api.NewPodMetrics(pod, string(event.Type), string(pod.Status.Phase))

		//log.Printf("%s SchedulingTime %dms, PrepTime %dms, Running Time %dms\n",
		//	podMetric.PodName, podMetric.GetSchedulingTime(), podMetric.GetPrepTime(), podMetric.GetRunningTime())
		//TODO introduce annotations and search by function annotation
		podNameNormalized := strings.Split(pod.ObjectMeta.Name, "-000")[0]
		switch {
		case event.Type == watch.Modified && len(pod.Status.HostIP) > 0:
			api.PodsMap.Store(podNameNormalized, api.Pod{
				PodName:     podNameNormalized,
				NodeName:    pod.Spec.NodeName,
				NodeIP:      pod.Status.HostIP,
				PodIp:       pod.Status.PodIP,
				Annotations: pod.ObjectMeta.Annotations,
			})
		case event.Type == watch.Deleted || event.Type == watch.Error:
			if podMetric.Exists() {
				delete(api.PodMetricsMap, podNameNormalized)
			}
			api.PodsMap.LoadAndDelete(podNameNormalized)
		}
	}
}
