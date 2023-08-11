package main

import (
	"os"

	mqttv1 "github.com/tomkiin/mosquitto-operator/api/v1"
	"github.com/tomkiin/mosquitto-operator/controller"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mqttv1.AddToScheme(scheme))
}

func main() {
	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
		o.TimeEncoder = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
	})
	if err != nil {
		ctrl.Log.Error(err, "Create Manager failed")
		os.Exit(1)
	}

	controller := &controller.MosquittoReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	if err := controller.Setup(mgr); err != nil {
		ctrl.Log.Error(err, "Start Controller failed")
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "Start Manager failed")
		os.Exit(1)
	}
}
