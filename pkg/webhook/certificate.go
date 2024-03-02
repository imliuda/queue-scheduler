package webhook

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	v1 "k8s.io/api/certificates/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/certificate"
	"k8s.io/kubernetes/pkg/controller/certificates/approver"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"os"
)

var (
	oidExtensionSubjectAltName = []int{2, 5, 29, 17}
)

func NewCertificateManager(kubeClient kubernetes.Interface) (certificate.Manager, error) {
	var clientsetFn certificate.ClientsetFunc
	if kubeClient != nil {
		clientsetFn = func(current *tls.Certificate) (kubernetes.Interface, error) {
			return kubeClient, nil
		}
	}
	certificateStore, err := certificate.NewFileStore(
		"queue-scheduler",
		"/tmp/queue-scheduler-cert",
		"/tmp/queue-scheduler-cert",
		"tls.crt",
		"tls.key")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize certificate store: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	getTemplate := func() *x509.CertificateRequest {
		return &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:   fmt.Sprintf("system:node:%s", hostname),
				Organization: []string{"system:nodes"},
			},
			DNSNames: []string{"queue-scheduler.kube-system.svc"},
			Extensions: []pkix.Extension{
				{
					Id:       oidExtensionSubjectAltName,
					Critical: false,
					Value:    []byte("DNS:queue-scheduler.kube-system.svc"),
				},
			},
		}
	}

	m, err := certificate.NewManager(&certificate.Config{
		ClientsetFn:      clientsetFn,
		GetTemplate:      getTemplate,
		SignerName:       "kubernetes.io/kubelet-serving",
		Usages:           []v1.KeyUsage{v1.UsageDigitalSignature, v1.UsageKeyEncipherment, v1.UsageServerAuth},
		CertificateStore: certificateStore,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server certificate manager: %v", err)
	}

	return m, nil
}

type csrapprover struct {
	f        framework.Handle
	informer cache.SharedIndexInformer
}

func NewCertificateApprover(ctx context.Context, f framework.Handle) (csrapprover, error) {
	s := csrapprover{
		informer: f.SharedInformerFactory().Certificates().V1().CertificateSigningRequests().Informer(),
	}

	return s, nil
}

func (s *csrapprover) OnAdd(obj interface{}, isInInitialList bool) {

	//TODO implement me
	panic("implement me")
}

func (s *csrapprover) OnUpdate(oldObj, newObj interface{}) {
	//TODO implement me
	panic("implement me")
}

func (s *csrapprover) OnDelete(obj interface{}) {
	//TODO implement me
	panic("implement me")
}

func (s *csrapprover) Start(ctx context.Context) {
	a := approver.NewCSRApprovingController(ctx, s.f.ClientSet(), s.f.SharedInformerFactory().Certificates().V1().CertificateSigningRequests())
	a.Run(ctx, 5)
}
