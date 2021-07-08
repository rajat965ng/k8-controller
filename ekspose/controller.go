package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	v12 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type controller struct {
	clientset      kubernetes.Interface
	deplister      v1.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, informer v12.DeploymentInformer) *controller {
	c := &controller{
		clientset:      clientset,
		deplister:      informer.Lister(),
		depCacheSynced: informer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addHandler,
			DeleteFunc: c.deleteHandler,
		},
	)

	return c
}

func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("starting controller !!")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Println("Error inside cache to be synced !!")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool  {
	item,shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	key,err := cache.MetaNamespaceKeyFunc(item)
	if err!=nil  {
		panic(err)
	}

	ns,name,err := cache.SplitMetaNamespaceKey(key)
	if err!=nil {
		panic(err)
		return false
	}

	err = c.syncDeployment(ns,name)
	if err!=nil {
		panic(err)
		return false
	}

	return true
}

func (c *controller) syncDeployment(ns string,name string) error  {
//create  service
//create ingress

	return nil
}

func (c *controller) addHandler(obj interface{}) {
	fmt.Println("Add was called ")
	c.queue.Add(obj)
}

func (c *controller) deleteHandler(obj interface{}) {
	fmt.Println("Delete was called ")
	c.queue.Add(obj)
}
