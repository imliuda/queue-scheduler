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

package webhook

import (
	"context"
	"github.com/imliuda/queue-scheduler/api/config"
	"github.com/imliuda/queue-scheduler/api/scheduling/v1alpha1"
	"github.com/imliuda/queue-scheduler/pkg/cache"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:admissionReviewVersions=v1alpha1,path=/mutate-v1alpha1-queues.scheduling.queue-scheduler.imliuda.github.com,mutating=true,failurePolicy=fail,groups="scheduling.queue-scheduler.imliuda.github.com",resources=queues,verbs=create;update,versions=v1alpha1,name=mpod.kb.io,sideEffects=None
// +kubebuilder:webhook:admissionReviewVersions=v1alpha1,path=/valiadate-v1alpha1-queues.scheduling.queue-scheduler.imliuda.github.com,mutating=false,failurePolicy=fail,groups="scheduling.queue-scheduler.imliuda.github.com",resources=queues,verbs=create;update;delete,versions=v1alpha1,name=mpod.kb.io,sideEffects=None

type QueueWebhook struct {
	client.Client
}

//
//func (w *QueueWebhook) validateQueue(ctx context.Context, q *v1alpha1.QueueConfig) (warnings admission.Warnings, err error) {
//	// validate min/max
//	// if a resource in max is unspecified, then it can use as much as possible
//	if !utils.LessOrEqual(q.Spec.Min, q.Spec.Max, true) {
//		return nil, fmt.Errorf("min resource of queue [%s] is larger then max", q.ObjectMeta.Name)
//	}
//
//	// validate parent queue
//	var parentName, name string
//	parts := strings.Split(name, ".")
//	if len(parts) > 1 {
//		parentName = strings.Join(parts[:len(parts)-1], ".")
//		name = parts[len(parts)-1]
//	}
//
//	if parentName != "" {
//		parentQueue := &v1alpha1.QueueConfig{}
//		err = w.Client.Get(ctx, types.NamespacedName{Namespace: "", Name: parentName}, parentQueue)
//		if err != nil {
//			return nil, fmt.Errorf("get parent queue [%s] error", parentName)
//		}
//
//		// validate max is smaller then parent
//		if !utils.LessOrEqual(q.Spec.Max, parentQueue.Spec.Max, true) {
//			return nil, fmt.Errorf("max resource of queue [%s] is larger then parent", name)
//		}
//
//		siblingQueues := &v1alpha1.QueueConfigList{}
//		if err = w.Client.List(ctx, siblingQueues, &client.ListOptions{
//			LabelSelector: labels.SelectorFromSet(labels.Set{"parent": parentName}),
//		}); err != nil {
//			return nil, fmt.Errorf("get sibling queues of queue [%s] error", name)
//		}
//
//		totalMin := corev1.ResourceList{}
//		for _, sq := range siblingQueues.Items {
//			if sq.Spec.Min != nil {
//				for resourceName := range sq.Spec.Min {
//					var quantity resource.Quantity
//					if v, ok := totalMin[resourceName]; ok {
//						quantity = v
//					}
//					quantity.Add(sq.Spec.Min[resourceName])
//					totalMin[resourceName] = quantity
//				}
//			}
//		}
//		if !utils.LessOrEqual(totalMin, parentQueue.Spec.Min, false) {
//			return nil, fmt.Errorf("sum of sbilings queue's min resources is greater then parent queue [%s] min resource", parentName)
//		}
//	}
//	return nil, nil
//}
//
//func (w *QueueWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
//	q := obj.(*v1alpha1.QueueConfig)
//
//	if q.ObjectMeta.Labels == nil || q.Labels["parent"] != q.ObjectMeta.Name {
//		return nil, fmt.Errorf("queue [%s] has empty labels, can check parent queue", q.ObjectMeta.Name)
//	}
//
//	if _, ok := q.ObjectMeta.Labels["parent"]; !ok {
//		return nil, fmt.Errorf("queue [%s] dos't have parent label", q.ObjectMeta.Name)
//	}
//
//	return w.validateQueue(ctx, q)
//}
//
//func (w *QueueWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
//	newQ := oldObj.(*v1alpha1.QueueConfig)
//
//	if newQ.Labels == nil || newQ.Labels["parent"] != newQ.ObjectMeta.Name {
//		return nil, fmt.Errorf("can't modify parent label")
//	}
//
//	return w.validateQueue(ctx, newQ)
//}
//
//func (w *QueueWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
//	q := obj.(*v1alpha1.QueueConfig)
//
//	childQueues := &v1alpha1.QueueConfigList{}
//	if err = w.Client.List(ctx, childQueues, &client.ListOptions{
//		LabelSelector: labels.SelectorFromSet(labels.Set{"parent": q.ObjectMeta.Name}),
//	}); err != nil {
//		return nil, fmt.Errorf("get child queues of queue [%s] error", q.ObjectMeta.Name)
//	}
//
//	if len(childQueues.Items) > 0 {
//		return nil, fmt.Errorf("can't delete queue [%s], because it has child queues", q.ObjectMeta.Name)
//	}
//	return nil, nil
//}
//
//func (w *QueueWebhook) Default(ctx context.Context, obj runtime.Object) error {
//	q := obj.(*v1alpha1.QueueConfig)
//
//	if q.Spec.Weight == 0 {
//		q.Spec.Weight = 1
//	}
//
//	if q.ObjectMeta.Labels == nil {
//		q.ObjectMeta.Labels = make(map[string]string)
//	}
//
//	parts := strings.Split(q.ObjectMeta.Name, ".")
//	if len(parts) > 1 {
//		q.ObjectMeta.Labels["parent"] = strings.Join(parts[:len(parts)-1], ".")
//	} else {
//		q.ObjectMeta.Labels["parent"] = ""
//	}
//
//	return nil
//}

type QueueConfigValidator struct {
}

func (q QueueConfigValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	newQ := obj.(*v1alpha1.QueueConfig)
	root := cache.FromConfig(newQ)
	if root.Validate() != nil {
		return nil, err
	}
	return nil, nil
}

func (q QueueConfigValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	newQ := newObj.(*v1alpha1.QueueConfig)
	root := cache.FromConfig(newQ)
	if root.Validate() != nil {
		return nil, err
	}
	return nil, nil
}

func (q QueueConfigValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// TODO make config configurable

func Setup(ctx context.Context, schema *runtime.Scheme, args *config.HierarchyQueueArgs, f framework.Handle) error {
	// setup cert signer
	//approver, err := NewCertificateApprover(ctx, f)
	//if err != nil {
	//	return err
	//}
	//approver.Start(ctx)

	// setup cert manager
	//selfSubjectAccessReview := &v1.SelfSubjectAccessReview{
	//	Spec: v1.SelfSubjectAccessReviewSpec{
	//		ResourceAttributes: &v1.ResourceAttributes{
	//			Namespace: "kube-system",
	//			Verb:      "create",
	//			Group:     "certificates.k8s.io",
	//			Version:   "*",
	//			Resource:  "certificatesigningrequests",
	//		},
	//	},
	//}
	//review, err := f.ClientSet().AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, selfSubjectAccessReview, metav1.CreateOptions{})
	//if err != nil {
	//	return errors.Wrap(err, "check create CSR permission failed")
	//}
	//if !review.Status.Allowed {
	//	return errors.NewRoot("no permission to create CSR")
	//}
	//
	//certManager, err := NewCertificateManager(f.ClientSet())
	//if err != nil {
	//	return err
	//}
	//certManager.Start()

	// setup webhook server
	//tlsOpt := func(c *tls.Config) {
	//	c.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	//		return certManager.Current(), nil
	//	}
	//}
	webhookServer := webhook.NewServer(webhook.Options{
		//TLSOpts: []func(*tls.Config){tlsOpt},
		Port:    args.WebHook.Port,
		CertDir: "/tmp/queue-scheduler-cert",
	})

	validator := &QueueConfigValidator{}
	webhookServer.Register("/valiadate-v1alpha1-queues.scheduling.queue-scheduler.imliuda.github.com",
		admission.WithCustomValidator(schema, &v1alpha1.QueueConfig{}, validator))

	go func() {
		if err := webhookServer.Start(ctx); err != nil {
			logger := klog.FromContext(ctx)
			logger.Error(nil, "Webhook server stopped")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}
	}()
	return nil
}
