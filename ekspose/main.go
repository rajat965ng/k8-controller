package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/Users/rajnigam/.kube/config", "location of kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error in loading kubeconfig %s\n",err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error in loading clientset %s\n",err.Error())
	}

	ch := make(chan struct{})
	informer := informers.NewSharedInformerFactory(clientset, 30*time.Minute)

	controller := newController(clientset, informer.Apps().V1().Deployments())

	informer.Start(ch)
	controller.run(ch)
	fmt.Println(informer)
}
