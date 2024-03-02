/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"context"
	"fmt"
	"github.com/imliuda/queue-scheduler/api/config"
	_ "github.com/imliuda/queue-scheduler/api/config/scheme"
	schedulingv1alpha1 "github.com/imliuda/queue-scheduler/api/scheduling/v1alpha1"
	schedulingcache "github.com/imliuda/queue-scheduler/pkg/cache"
	"github.com/imliuda/queue-scheduler/pkg/queue"
	"github.com/imliuda/queue-scheduler/pkg/webhook"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(schedulingv1alpha1.AddToScheme(scheme))
}

var _ framework.PreFilterPlugin = &HierarchyQueue{}

type HierarchyQueue struct {
	handle      framework.Handle
	queueStates *queue.States
}

func New(ctx context.Context, configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	args, ok := configuration.(*config.HierarchyQueueArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type HierarchyQueueArgs, got %T", configuration)
	}

	if err := webhook.Setup(ctx, scheme, args, f); err != nil {
		return nil, errors.Wrap(err, "setup HierarchyQueue webhook error")
	}

	cache := schedulingcache.NewCache(args, f)
	cache.Run(ctx)

	return &HierarchyQueue{
		handle: f,
	}, nil
}

func (s HierarchyQueue) PreFilter(ctx context.Context, state *framework.CycleState, p *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	//TODO implement me
	panic("implement me")
}

func (s HierarchyQueue) PreFilterExtensions() framework.PreFilterExtensions {
	//TODO implement me
	panic("implement me")
}

func (s HierarchyQueue) Name() string {
	return "HierarchyQueue"
}
