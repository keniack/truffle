package metrics

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"polaris/truffle/pkg/common"
)

func StartPodMetrics() {

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
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// initial list
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}

	// watch future changes to Pods
	watcher, err := clientset.CoreV1().Pods("").Watch(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()
	if common.Debug {
		common.DebugLog.Printf("--- POD Watch (namespace %v) ----\n", "")
	}
	for event := range ch {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			common.DebugLog.Fatal("unexpected type")
		}
		if common.Debug {
			common.DebugLog.Printf("Event [%s] PodName [%s] Status [%s] NodeName [%s] HostIP [%s] PodIp [%s]\n",
				event.Type, pod.ObjectMeta.Name, pod.Status.Phase, pod.Spec.NodeName, pod.Status.HostIP, pod.Status.PodIPs)
		}

		var podMetric, _ = NewPodMetrics(pod, string(event.Type), string(pod.Status.Phase))
		if common.Debug {
			common.DebugLog.Printf("%s SchedulingTime %dms, PrepTime %dms, Running Time %dms\n",
				podMetric.PodName, podMetric.GetSchedulingTime(), podMetric.GetPrepTime(), podMetric.GetRunningTime())
		}
	}
}
