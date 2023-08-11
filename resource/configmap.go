package resource

import (
	"fmt"

	mqttv1 "github.com/tomkiin/mosquitto-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigMapName() string {
	return "mosquitto-config"
}

func CreateConfigMap(ins *mqttv1.Mosquitto) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName(),
			Namespace: ins.Namespace,
		},
		Data: map[string]string{
			"mosquitto-node.conf":   ins.Spec.Conf,
			"mosquitto-master.conf": createMasterConf(ins),
		},
	}
}

func IsConfigMapChanged(cfg, oldCfg *corev1.ConfigMap) bool {
	if cfg.Data == nil && oldCfg.Data == nil {
		return false
	}

	if cfg.Data == nil || oldCfg.Data == nil {
		return true
	}

	return cfg.Data["mosquitto-node.conf"] != oldCfg.Data["mosquitto-node.conf"] ||
		cfg.Data["mosquitto-master.conf"] != oldCfg.Data["mosquitto-master.conf"]
}

func createMasterConf(ins *mqttv1.Mosquitto) string {
	conf := ins.Spec.Conf
	for i := 1; i < int(ins.Spec.Count); i++ {
		pod := fmt.Sprintf("%s-%d", StatefulSetName(), i)
		conf += fmt.Sprintf("\nconnection %s\naddress %s.%s.%s.svc.cluster.local\ntopic # both 0\n", pod, pod, HeadlessServiceName(), ins.Namespace)
	}

	return conf
}
