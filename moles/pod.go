// Sample K8s pod event observer. Watches pods coming up or down
package moles

import (
	"log"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

const resyncPeriod = 30 * time.Minute

type callbacks struct {
	upFn func(*api.Pod)
	dnFn func(*api.Pod)
}

// Watch pod activity on the cluster and triggers custom behavior
func Watch(masterURL string, up, down func(*api.Pod)) *client.Client {
	config := &restclient.Config{
		Host:     masterURL,
		Insecure: true,
	}

	master, err := client.New(config)
	if err != nil {
		log.Fatalln("K8s master connect failed!", err)
	}

	watchPods(master, callbacks{upFn: up, dnFn: down})

	return master
}

func watchPods(client *client.Client, cb callbacks) {
	watchlist := cache.NewListWatchFromClient(
		client,
		"pods",
		api.NamespaceDefault,
		fields.Everything(),
	)

	_, eController := framework.NewInformer(
		watchlist,
		&api.Pod{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    cb.up,
			DeleteFunc: cb.down,
		},
	)

	go eController.Run(wait.NeverStop)
}

func (cb callbacks) up(obj interface{}) {
	pod := obj.(*api.Pod)
	log.Printf("Pod created: %s (%s)\n", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)

	cb.upFn(pod)
}

func (cb callbacks) down(obj interface{}) {
	pod := obj.(*api.Pod)
	log.Printf("Pod deleted: %s (%s)\n", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)

	cb.dnFn(pod)
}
