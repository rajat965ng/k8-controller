package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)


func main()  {
	kubeconfig := flag.String("kubeconfig","/Users/rajnigam/.kube/config","location of kubeconfig file")
	config,err  := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	config.Timeout = 120 * time.Second
	if err!=nil {
		fmt.Printf("The error is in Config: %s",err.Error())
		fmt.Println("Applying Incluster Config !!!")
		config,err = rest.InClusterConfig()
		if err!=nil {
			fmt.Printf("The error is in Incluster Config: %s",err.Error())
		}
	}

	//runtime.Object()
	clientset,err := kubernetes.NewForConfig(config)
	if err!= nil {
		fmt.Printf("The error is in ClientSet: %s",err.Error())
	}

	pods,err := clientset.CoreV1().Pods("default").List(context.Background(),v1.ListOptions{})

	if err!=nil {
		fmt.Printf("The error is in while reading pods: %s",err.Error())
	}

	fmt.Println("pods from default namespace")
	for _,pod := range pods.Items {
		fmt.Printf("\n pod name: %s",pod.Name)
	}
	fmt.Println("\n ------------ ")


	deployments,err := clientset.AppsV1().Deployments("default").List(context.Background(),v1.ListOptions{})

	if err!=nil {
		fmt.Printf("The error is in while reading deployments: %s",err.Error())
	}

	fmt.Println("deployments from default namespace")
	for _,deployment := range deployments.Items {
		fmt.Printf("\n deployment: %s",deployment.Name)
	}
	fmt.Println("\n -----------")
}