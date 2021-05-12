package main

import (
	"flag"
	"os"

	"github.com/octohelm/cuemod/internal/version"
	"github.com/octohelm/cuemod/pkg/cuemoperator"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func start(ctrlOpt ctrl.Options) error {
	restConfig := ctrl.GetConfigOrDie()

	mgr, err := ctrl.NewManager(restConfig, ctrlOpt)
	if err != nil {
		return errors.Wrap(err, "unable to start manager")
	}

	if err := cuemoperator.SetupWithManager(mgr); err != nil {
		return errors.Wrap(err, "unable to create controller")
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return errors.Wrap(err, "problem running manager")
	}
	return nil
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	ctrlOpt := ctrl.Options{
		Scheme:           scheme,
		Port:             9443,
		LeaderElectionID: "asd1231ax.octohelm.tech",
		Logger:           ctrl.Log.WithValues("cuem-operator", version.Version),
	}

	flag.StringVar(&ctrlOpt.Namespace, "watch-namespace", os.Getenv("WATCH_NAMESPACE"), "watch namespace")
	flag.StringVar(&ctrlOpt.MetricsBindAddress, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&ctrlOpt.LeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.Parse()

	if err := start(ctrlOpt); err != nil {
		ctrl.Log.WithName("setup").Error(err, "")
	}
}
