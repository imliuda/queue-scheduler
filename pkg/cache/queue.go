package cache

import (
	"fmt"
	"github.com/imliuda/queue-scheduler/api/scheduling/v1alpha1"
	"github.com/imliuda/queue-scheduler/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
	"strings"
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

func NewRoot() *Queue {
	return &Queue{
		Name:       "root",
		FullName:   "root",
		Properties: make(map[string]string),
		Queues:     make([]*Queue, 0),
		Parent:     nil,
	}
}

func FromConfig(queue *v1alpha1.QueueConfig) *Queue {
	queue = queue.DeepCopy()

	root := NewRoot()
	root.Update(queue)

	return root
}

func update(root *Queue, target v1alpha1.Queue) {
	currentChildrenMap := make(map[string]*Queue)
	targetChildrenMap := make(map[string]v1alpha1.Queue)
	childrenToDelete := make(map[string]bool)
	childrenToAdd := make([]v1alpha1.Queue, 0)
	childrenToUpdate := make(map[string]bool)

	// update self first
	root.Min = target.Min
	root.Max = target.Max
	root.Weight = target.Weight
	if target.Properties != nil {
		root.Properties = target.Properties
	}

	for _, q := range root.Queues {
		currentChildrenMap[q.Name] = q
	}
	for _, q := range target.Queues {
		targetChildrenMap[q.Name] = q
	}

	for n, q := range targetChildrenMap {
		if _, ok := currentChildrenMap[n]; !ok {
			childrenToAdd = append(childrenToAdd, q)
		}
	}

	for n, _ := range currentChildrenMap {
		if _, ok := targetChildrenMap[n]; ok {
			childrenToUpdate[n] = true
		} else {
			childrenToDelete[n] = true
		}
	}

	if len(childrenToUpdate) > 0 {
		for n, _ := range childrenToUpdate {
			update(currentChildrenMap[n], targetChildrenMap[n])
		}
	}

	if len(childrenToAdd) > 0 {
		for _, child := range childrenToAdd {
			q := &Queue{
				Name:       child.Name,
				FullName:   fmt.Sprintf("%s.%s", root.FullName, child.Name),
				Queues:     make([]*Queue, 0),
				Properties: make(map[string]string),
				Parent:     root,
			}
			update(q, child)
			root.Queues = append(root.Queues, q)
		}
	}

	if len(childrenToDelete) > 0 {
		children := make([]*Queue, 0)
		for _, q := range root.Queues {
			if !childrenToDelete[q.Name] {
				children = append(children, q)
			}
		}
		root.Queues = children
	}
}

func (q *Queue) Update(queue *v1alpha1.QueueConfig) {
	queue = queue.DeepCopy()

	target := v1alpha1.Queue{
		Name:       "root",
		Queues:     queue.Queues,
		Properties: make(map[string]string),
	}

	if len(queue.Queues) == 1 && queue.Queues[0].Name == "root" {
		target = queue.Queues[0]
	}

	update(q, target)

	klog.V(5).Info("Queue updated\n", q.Dumps())
}

func dump(q *Queue, sb *strings.Builder, indent int) {
	sb.WriteString(fmt.Sprintf("%s %s Min: %v Max %v Properties %v\n",
		strings.Repeat(" ", indent), q.Name, q.Min, q.Max, q.Properties))
	for _, c := range q.Queues {
		dump(c, sb, indent+4)
	}
}

func (q *Queue) Dumps() string {
	sb := &strings.Builder{}
	dump(q, sb, 0)
	return sb.String()
}

func (q *Queue) Reset() {
	q.Min = v1.ResourceList{}
	q.Max = v1.ResourceList{}
	q.Properties = make(map[string]string)
	q.Queues = make([]*Queue, 0)

	klog.V(5).Info("Queue updated\n", q.Dumps())
}

func (q *Queue) UpdateState(state *States) {

}

func (q *Queue) SumQueuesMin(queues []*Queue) v1.ResourceList {
	totalMin := v1.ResourceList{}
	for _, sq := range queues {
		if sq.Min != nil {
			for resourceName := range sq.Min {
				var quantity resource.Quantity
				if v, ok := totalMin[resourceName]; ok {
					quantity = v
				}
				quantity.Add(sq.Min[resourceName])
				totalMin[resourceName] = quantity
			}
		}
	}
	return totalMin
}

func (q *Queue) GetQueuesMax(queues []*Queue) v1.ResourceList {
	max := v1.ResourceList{}
	for _, q := range queues {
		for resourceName, resourceValue := range q.Max {
			if v, ok := max[resourceName]; !ok || v.Cmp(resourceValue) <= 0 {
				max[resourceName] = resourceValue
			}
		}
	}
	return max
}

func (q *Queue) Validate() error {
	klog.V(5).Info("validating Queue create request", "q", q)

	// validate self min less or equal max
	if !utils.LessEqual(q.Min, q.Max, true) {
		return fmt.Errorf("min resource of Queue [%s] is larger then max", q.FullName)
	}

	if q.Parent != nil {
		// validate max is smaller then parent max
		parent := q.Parent
		if !utils.LessEqual(q.Max, parent.Max, true) {
			return fmt.Errorf("max resource of Queue [%s] is larger then parent", q.FullName)
		}

		// validate sum(min of siblings) is not greater than parent's min
		siblings := parent.Queues
		totalMin := q.SumQueuesMin(siblings)

		if !utils.LessEqual(totalMin, parent.Min, false) {
			return fmt.Errorf("sum of sbilings Queue's min resources is greater then parent Queue [%s] min", q.Parent.FullName)
		}
	}

	// validate sum(min of child) is not greater than current Queue min
	totalMin := q.SumQueuesMin(q.Queues)

	if !utils.LessEqual(totalMin, q.Min, false) {
		return fmt.Errorf("sum of sbilings Queue's min resources is greater then current Queue [%s] min", q.FullName)
	}

	// validate max(child) is not greater than current Queue max
	qmax := q.GetQueuesMax(q.Queues)
	if !utils.LessEqual(qmax, q.Max, false) {
		return fmt.Errorf("one of children's max is greater than current Queue [%s] max", q.FullName)
	}

	return nil
}
