package cache

import (
	"context"
	"github.com/imliuda/queue-scheduler/api/config"
	"github.com/imliuda/queue-scheduler/api/scheduling/v1alpha1"
	"github.com/imliuda/queue-scheduler/pkg/generated/clientset/versioned"
	"github.com/imliuda/queue-scheduler/pkg/generated/informers/externalversions"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sync"
	"time"
)

type Cache struct {
	queueInformer cache.SharedIndexInformer
	queue         *Queue
	mu            sync.Mutex
}

func NewCache(args *config.HierarchyQueueArgs, f framework.Handle) *Cache {
	cache := &Cache{}

	clientSet := versioned.NewForConfigOrDie(f.KubeConfig())
	externalSharedInformerFactory := externalversions.NewSharedInformerFactory(clientSet, 10*time.Minute)
	informer := externalSharedInformerFactory.Scheduling().V1alpha1().QueueConfigs().Informer()
	informer.AddEventHandler(cache)

	cache.queueInformer = informer
	cache.queue = NewRoot()

	return cache
}

func (c *Cache) Run(ctx context.Context) {
	c.queueInformer.Run(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), c.queueInformer.HasSynced)
}

func (c *Cache) OnAdd(obj interface{}, isInInitialList bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch obj.(type) {
	case *v1alpha1.QueueConfig:
		q := obj.(*v1alpha1.QueueConfig)
		if q.Name == "queue-scheduler" {
			c.queue.Update(q)
		} else {
			klog.Info("QueueConfig config is not named queue-scheduler, will not create...")
		}
	}
}

func (c *Cache) OnUpdate(oldObj, newObj interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch newObj.(type) {
	case *v1alpha1.QueueConfig:
		q := newObj.(*v1alpha1.QueueConfig)
		if q.Name == "queue-scheduler" {
			c.queue.Update(q)
		} else {
			klog.Info("QueueConfig config is not named queue-scheduler, will not update...")
		}
	}
}

func (c *Cache) OnDelete(obj interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	q := obj.(*v1alpha1.QueueConfig)
	if q.Name == "queue-scheduler" {
		c.queue.Reset()
	} else {
		klog.Info("QueueConfig config is not named queue-scheduler, will not update...")
	}
}

func (c *Cache) Snapshot() {
	c.mu.Lock()
	defer c.mu.Unlock()

}
