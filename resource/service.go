package resource

import (
	mqttv1 "github.com/tomkiin/mosquitto-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func HeadlessServiceName() string {
	return "mosquitto-headless-svc"
}

func CreateHeadlessService(ins *mqttv1.Mosquitto) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadlessServiceName(),
			Namespace: ins.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: StatefulSetLabels(),
			Ports: []corev1.ServicePort{
				{
					Name:       "mqtt-port",
					Protocol:   corev1.ProtocolTCP,
					Port:       1883,
					TargetPort: intstr.FromInt(1883),
				},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}
}
