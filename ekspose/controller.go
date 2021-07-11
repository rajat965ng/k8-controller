package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwk1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	// If everything goes fine with syncDeployment then "item" should not be processed again
	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("Error in getting MetaNamespaceKeyFunc %s\n",err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("Error in getting SplitMetaNamespaceKey %s\n",err.Error())
		return false
	}

	// Check if  the object has been deleted from k8 cluster
	ctx := context.Background()
	_, err = c.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("Deployment %s was deleted !!\n", name)
		err = c.clientset.CoreV1().Services(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Error deleting %s svc, error %s\n", name, err.Error())
			return false
		}

		err = c.clientset.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Error deleting %s ingress, error %s\n", name, err.Error())
			return false
		}
		return true
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		//re-try logic
		fmt.Printf("\nError during sync deployments: %s \n", err.Error())
		return false
	}

	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()
	dep, err := c.deplister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf("Error in listing deployment: %s \n",err.Error())
	}

	//create  service
	//we have to modify this to  figure out  the port
	//our deployment container is listening on
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabels(*dep),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}

	//initialize  svc  for creating  ingress
	s, err := c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("\nError in creating service: %s", err.Error())
	}
	//create ingress
	return createIngress(ctx, c.clientset, *s)
}

func createIngress(ctx context.Context, client kubernetes.Interface, svc corev1.Service) error {
	pathType := "Prefix"
	ingress := nwk1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: nwk1.IngressSpec{
			Rules: []nwk1.IngressRule{
				{
					IngressRuleValue: nwk1.IngressRuleValue{
						HTTP: &nwk1.HTTPIngressRuleValue{
							Paths: []nwk1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/%s", svc.Name),
									PathType: (*nwk1.PathType)(&pathType),
									Backend: nwk1.IngressBackend{
										Service: &nwk1.IngressServiceBackend{
											Name: svc.Name,
											Port: nwk1.ServiceBackendPort{
												Number: 80,
											},
										},
									}},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.NetworkingV1().Ingresses(svc.Namespace).Create(ctx, &ingress, metav1.CreateOptions{})
	return err
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}

func (c *controller) addHandler(obj interface{}) {
	fmt.Println("Add was called ")
	c.queue.Add(obj)
}

func (c *controller) deleteHandler(obj interface{}) {
	fmt.Println("Delete was called ")
	c.queue.Add(obj)
}
