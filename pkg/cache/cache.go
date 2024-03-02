package cache

import (
	"context"
	"github.com/imliuda/queue-scheduler/api/config"
	"github.com/imliuda/queue-scheduler/pkg/generated/clientset/versioned"
	"github.com/imliuda/queue-scheduler/pkg/generated/informers/externalversions"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sync"
	"time"
)

type Cache struct {
	queueInformer cache.SharedIndexInformer
	mu            sync.Mutex
}

func NewCache(args *config.HierarchyQueueArgs, f framework.Handle) *Cache {
	cache := &Cache{}

	clientSet := versioned.NewForConfigOrDie(f.KubeConfig())
	externalSharedInformerFactory := externalversions.NewSharedInformerFactory(clientSet, 10*time.Minute)
	informer := externalSharedInformerFactory.Scheduling().V1alpha1().Queues().Informer()
	informer.AddEventHandler(cache)
	cache.queueInformer = informer

	return cache
}

func (c *Cache) Run(ctx context.Context) {
	c.queueInformer.Run(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), c.queueInformer.HasSynced)
}

func (c *Cache) OnAdd(obj interface{}, isInInitialList bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

}

func (c *Cache) OnUpdate(oldObj, newObj interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *Cache) OnDelete(obj interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *Cache) Snapshot() {
	c.mu.Lock()
	defer c.mu.Unlock()
}
