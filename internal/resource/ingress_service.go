package resource

import (
	"fmt"

	rabbitmqv1beta1 "github.com/pivotal/rabbitmq-for-kubernetes/api/v1beta1"
	"github.com/pivotal/rabbitmq-for-kubernetes/internal/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (builder *RabbitmqResourceBuilder) IngressService() *IngressServiceBuilder {
	return &IngressServiceBuilder{
		Instance: builder.Instance,
		Scheme:   builder.Scheme,
	}
}

type IngressServiceBuilder struct {
	Instance *rabbitmqv1beta1.RabbitmqCluster
	Scheme   *runtime.Scheme
}

func (builder *IngressServiceBuilder) Build() (runtime.Object, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      builder.Instance.ChildResourceName("ingress"),
			Namespace: builder.Instance.Namespace,
			Labels:    metadata.GetLabels(builder.Instance.Name, builder.Instance.Labels),
		},
		Spec: corev1.ServiceSpec{
			Selector: metadata.LabelSelector(builder.Instance.Name),
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     5672,
					Name:     "amqp",
				},
				{
					Protocol: corev1.ProtocolTCP,
					Port:     15672,
					Name:     "http",
				},
				{
					Protocol: corev1.ProtocolTCP,
					Port:     15692,
					Name:     "prometheus",
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(builder.Instance, service, builder.Scheme); err != nil {
		return nil, fmt.Errorf("failed setting controller reference: %v", err)
	}

	builder.setServiceType(service)
	builder.setAnnotations(service)

	return service, nil
}

func (builder *IngressServiceBuilder) setServiceType(service *corev1.Service) {
	var serviceType = "ClusterIP"
	if builder.Instance.Spec.Service.Type != "" {
		serviceType = builder.Instance.Spec.Service.Type
	}
	service.Spec.Type = corev1.ServiceType(serviceType)
}

func (builder *IngressServiceBuilder) Update(object runtime.Object) error {
	service := object.(*corev1.Service)
	builder.setAnnotations(service)
	service.Labels = metadata.GetLabels(builder.Instance.Name, builder.Instance.Labels)
	return nil
}

func (builder *IngressServiceBuilder) setAnnotations(service *corev1.Service) {
	if builder.Instance.Spec.Service.Annotations != nil {
		service.Annotations = metadata.ReconcileAnnotations(service.Annotations, builder.Instance.Annotations, builder.Instance.Spec.Service.Annotations)
	} else {
		service.Annotations = metadata.ReconcileAnnotations(service.Annotations, builder.Instance.Annotations)
	}
}
