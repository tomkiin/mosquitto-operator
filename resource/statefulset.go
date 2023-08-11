package resource

import (
	mqttv1 "github.com/tomkiin/mosquitto-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const shell = `set -ex
pod=$(hostname | awk -F"-" '{print $3}')
if [[ $pod -eq 0 ]]; then
	cp /opt/mosquitto/config/mosquitto-master.conf /mosquitto/config/mosquitto.conf
else
	cp /opt/mosquitto/config/mosquitto-node.conf /mosquitto/config/mosquitto.conf
fi;
`

func StatefulSetName() string {
	return "mosquitto-cluster"
}

func StatefulSetLabels() map[string]string {
	return map[string]string{"app": StatefulSetName()}
}

func CreateStatefulSet(ins *mqttv1.Mosquitto) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      StatefulSetName(),
			Namespace: ins.Namespace,
			Labels:    StatefulSetLabels(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &ins.Spec.Count,
			Selector: &metav1.LabelSelector{
				MatchLabels: StatefulSetLabels(),
			},
			ServiceName: HeadlessServiceName(),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: StatefulSetLabels(),
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ConfigMapName(),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "mosquitto-master.conf",
											Path: "mosquitto-master.conf",
										},
										{
											Key:  "mosquitto-node.conf",
											Path: "mosquitto-node.conf",
										},
									},
								},
							},
						},
						{
							Name: "empty-dir",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "load-conf",
							Image: ins.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "conf",
									MountPath: "/opt/mosquitto/config",
								},
								{
									Name:      "empty-dir",
									MountPath: "/mosquitto/config",
								},
							},
							Command: []string{
								"sh",
								"-c",
								shell,
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  StatefulSetName(),
							Image: ins.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "empty-dir",
									MountPath: "/mosquitto/config",
								},
							},
						},
					},
				},
			},
		},
	}
}
