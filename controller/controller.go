package controller

import (
	"context"

	"github.com/go-logr/logr"
	mqttv1 "github.com/tomkiin/mosquitto-operator/api/v1"
	"github.com/tomkiin/mosquitto-operator/resource"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ reconcile.Reconciler = (*MosquittoReconciler)(nil)

type MosquittoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (m *MosquittoReconciler) Setup(mgr ctrl.Manager) error {
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1.Mosquitto{}).
		Build(m)
	if err != nil {
		return err
	}

	watchObjs := []client.Object{
		&corev1.ConfigMap{},
		&corev1.Service{},
		&appsv1.StatefulSet{},
	}
	handler := &handler.EnqueueRequestForOwner{OwnerType: &mqttv1.Mosquitto{}, IsController: true}
	for _, obj := range watchObjs {
		if err := c.Watch(&source.Kind{Type: obj}, handler); err != nil {
			return err
		}
	}

	return nil
}

func (m *MosquittoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithName("MosquittoController")

	ins := &mqttv1.Mosquitto{}
	if err := m.Get(ctx, req.NamespacedName, ins); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	syncHandlers := []syncHandler{
		m.syncConfigMap,
		m.syncStatefulSet,
		m.syncHeadlessService,
	}
	for _, h := range syncHandlers {
		if err := h(ctx, ins, log); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

type syncHandler func(context.Context, *mqttv1.Mosquitto, logr.Logger) error

func (m *MosquittoReconciler) syncConfigMap(ctx context.Context, ins *mqttv1.Mosquitto, log logr.Logger) error {
	cfg := resource.CreateConfigMap(ins)
	if err := ctrl.SetControllerReference(ins, cfg, m.Scheme); err != nil {
		return err
	}

	oldCfg := &corev1.ConfigMap{}
	if err := m.Get(ctx, types.NamespacedName{Name: cfg.Name, Namespace: ins.Namespace}, oldCfg); err != nil {
		if errors.IsNotFound(err) {
			log.Info("create ConfigMap")
			return m.Create(ctx, cfg)
		}

		return err
	}

	if err := m.Update(ctx, cfg); err != nil {
		return err
	}

	if resource.IsConfigMapChanged(cfg, oldCfg) {
		log.Info("ConfigMap changed, reloading cluster")
		ins.Status.ClusterReloading = true
		return m.Status().Update(ctx, ins)
	}

	return nil
}

func (m *MosquittoReconciler) syncStatefulSet(ctx context.Context, ins *mqttv1.Mosquitto, log logr.Logger) error {
	sts := resource.CreateStatefulSet(ins)
	if err := ctrl.SetControllerReference(ins, sts, m.Scheme); err != nil {
		return err
	}

	if err := m.Get(ctx, types.NamespacedName{Name: sts.Name, Namespace: ins.Namespace}, &appsv1.StatefulSet{}); err != nil {
		if errors.IsNotFound(err) {
			log.Info("create StatefulSet")
			return m.Create(ctx, sts)
		}

		return err
	}

	if err := m.Update(ctx, sts); err != nil {
		return err
	}

	if ins.Status.ClusterReloading {
		log.Info("reloading cluster")
		if err := m.DeleteAllOf(ctx, &corev1.Pod{},
			client.InNamespace(ins.Namespace),
			client.MatchingLabels(resource.StatefulSetLabels()),
		); err != nil {
			return err
		}

		ins.Status.ClusterReloading = false
		return m.Status().Update(ctx, ins)
	}

	return nil
}

func (m *MosquittoReconciler) syncHeadlessService(ctx context.Context, ins *mqttv1.Mosquitto, log logr.Logger) error {
	svc := resource.CreateHeadlessService(ins)
	if err := ctrl.SetControllerReference(ins, svc, m.Scheme); err != nil {
		return err
	}

	if err := m.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: ins.Namespace}, &corev1.Service{}); err != nil {
		if errors.IsNotFound(err) {
			log.Info("create Headless Service")
			return m.Create(ctx, svc)
		}

		return err
	}

	return m.Update(ctx, svc)
}
