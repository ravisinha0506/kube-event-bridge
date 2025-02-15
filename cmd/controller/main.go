package main

import (
	"context"

	"github.com/ellistarn/kube-event-bridge/pkg/apis"
	"github.com/ellistarn/kube-event-bridge/pkg/controllers/event"
	"github.com/ellistarn/kube-event-bridge/pkg/controllers/eventrule"
	"github.com/ellistarn/kube-event-bridge/pkg/controllers/sqstarget"
	"github.com/go-logr/zapr"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	ctx := signals.NewContext()
	ctx = logging.WithLogger(ctx, lo.Must(zap.NewDevelopment()).Sugar())

	lo.Must0(apis.AddToScheme(scheme.Scheme))
	mgr := lo.Must(controllerruntime.NewManager(controllerruntime.GetConfigOrDie(), controllerruntime.Options{
		Logger: zapr.NewLogger(logging.FromContext(ctx).Desugar()),
		Scheme: scheme.Scheme,
	}))

	for _, register := range []func(context.Context, manager.Manager) error{
		event.Register,
		eventrule.Register,
		sqstarget.Register,
		slacktarget.Register,
	} {
		lo.Must0(register(ctx, mgr))
	}
	lo.Must0(mgr.Start(ctx))
}
