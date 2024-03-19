package main

import (
	"flag"
	_ "github.com/cossim/coss-server/docs"
	"github.com/cossim/coss-server/internal/msg/interface/grpc"
	"github.com/cossim/coss-server/internal/msg/interface/http"
	ctrl "github.com/cossim/coss-server/pkg/alias"
	"github.com/cossim/coss-server/pkg/config"
	"github.com/cossim/coss-server/pkg/discovery"
	"github.com/cossim/coss-server/pkg/healthz"
	"github.com/cossim/coss-server/pkg/manager/signals"
)

var (
	discover          bool
	register          bool
	remoteConfig      bool
	remoteConfigAddr  string
	remoteConfigToken string
	metricsAddr       string
	httpProbeAddr     string
	grpcProbeAddr     string
)

func init() {
	flag.BoolVar(&discover, "discover", false, "Enable service discovery")
	flag.BoolVar(&register, "register", false, "Enable service register")
	flag.BoolVar(&remoteConfig, "remote-config", false, "Load config from remote source")
	flag.StringVar(&remoteConfigAddr, "config-center-addr", "", "Address of the config center")
	flag.StringVar(&remoteConfigToken, "config-center-token", "", "Token for accessing the config center")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":9090", "The address the metric endpoint binds to")
	flag.StringVar(&httpProbeAddr, "http-health-probe-bind-address", ":9091", "The address to bind the http health probe endpoint")
	flag.StringVar(&grpcProbeAddr, "grpc-health-probe-bind-address", ":9092", "The address to bind the grpc health probe endpoint")
	flag.Parse()
}

func main() {
	svc := &grpc.Handler{}
	mgr, err := ctrl.NewManager(config.GetConfigOrDie(), ctrl.Options{
		Grpc: ctrl.GRPCServer{
			GRPCService:         svc,
			HealthzCheckAddress: grpcProbeAddr,
		},
		Http: ctrl.HTTPServer{
			HTTPService:        &http.Handler{MsgClient: svc},
			HealthCheckAddress: httpProbeAddr,
		},
		Config: ctrl.Config{
			LoadFromConfigCenter: remoteConfig,
			RemoteConfigAddr:     remoteConfigAddr,
			RemoteConfigToken:    remoteConfigToken,
			Hot:                  true,
			Key:                  "service/msg",
			Keys: []string{
				discovery.CommonMySQLConfigKey,
				discovery.CommonRedisConfigKey,
				discovery.CommonOssConfigKey,
				discovery.CommonMQConfigKey,
			},
			Registry: ctrl.Registry{
				Discover: discover,
				Register: register,
			},
		},
		MetricsBindAddress: metricsAddr,
	})
	if err != nil {
		panic(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(err)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		panic(err)
	}
}