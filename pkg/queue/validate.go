package queue

import (
	"fmt"
	"github.com/imliuda/queue-scheduler/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
)

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
