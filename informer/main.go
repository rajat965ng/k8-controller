package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main()  {
	kubeconfig := flag.String("kubeconfig","/Users/rajnigam/.kube/config","location of kubeconfig file")
	config,err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err!=nil {
		panic(err)
		config,err = rest.InClusterConfig()
		if err!=nil {
			panic(err)
		}
	}

	clientset,err := kubernetes.NewForConfig(config)
	if err!=nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset,30*time.Second)

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add was called !!")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update  was called !!")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete was called !!")
		},
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	pod,err := podInformer.Lister().Pods("default").Get("default")

	fmt.Println(pod)
	/*pods,err := clientset.CoreV1().Pods("default").List(context.Background(),v1.ListOptions{})
	if err!=nil {
		panic(err)
	}

	for _,pod := range pods.Items {
		fmt.Printf("The name of pod is: %s \n ",pod.Name)
	}*/
}