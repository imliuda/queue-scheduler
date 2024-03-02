package queue

import (
	"fmt"
	"github.com/imliuda/queue-scheduler/api/scheduling/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"sync"
)

var queueMu = sync.Mutex{}

type Queue struct {
	Name       string
	FullName   string
	Min        v1.ResourceList
	Max        v1.ResourceList
	Weight     int
	Properties map[string]string
	Queues     []*Queue
	Parent     *Queue
}

func New() *Queue {
	return &Queue{
		Name: "root",
	}
}

func addToRoot(root *Queue, children []v1alpha1.Queue) {
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		q := &Queue{
			Name:       child.Name,
			FullName:   fmt.Sprintf("%s.%s", root.FullName, child.Name),
			Min:        child.Min,
			Max:        child.Max,
			Weight:     child.Weight,
			Properties: child.Properties,
			Queues:     make([]*Queue, 0),
			Parent:     root,
		}

		addToRoot(q, child.Queues)

		root.Queues = append(root.Queues, q)
	}
}

func FromConfig(queue *v1alpha1.QueueConfig) *Queue {
	var root *Queue

	if len(queue.Queues) == 1 && queue.Queues[0].Name == "root" {
		cq := queue.Queues[0]
		root = &Queue{
			Name:       "root",
			FullName:   "root",
			Min:        cq.Min,
			Max:        cq.Max,
			Weight:     cq.Weight,
			Properties: cq.Properties,
			Queues:     make([]*Queue, 0),
			Parent:     nil,
		}
	} else {
		root = &Queue{
			Name:     "root",
			FullName: "root",
			Queues:   make([]*Queue, 0),
			Parent:   nil,
		}
	}

	addToRoot(root, queue.Queues)

	return root
}

func updateQueues(currents []*Queue, updates []v1alpha1.Queue) {
	//currentNameToQueues := make(map[string]*Queue)
	//updateNameToQueues := make(map[string]v1alpha1.Queue)
}

func (q *Queue) Update(queue *v1alpha1.QueueConfig) {
	queueMu.Lock()
	defer queueMu.Unlock()

	if len(queue.Queues) == 1 && queue.Queues[0].Name == "root" {
		cq := queue.Queues[0]
		q.Min = cq.Min
		q.Max = cq.Max
		q.Weight = cq.Weight
		q.Properties = cq.Properties
	} else {
		updateQueues(q.Queues, queue.Queues)
	}
}

func (q *Queue) UpdateState(state *States) {
	queueMu.Lock()
	defer queueMu.Unlock()

}
