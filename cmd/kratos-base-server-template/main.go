package main

import (
	"flag"
	"github.com/go-kratos/kratos/contrib/registry/zookeeper/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
	"os"
	"time"

	"kratos-base-template/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	//Name string
	Name = ""
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _       = os.Hostname()
	randomId, _ = uuid.NewUUID()
	instanceId  = id + "-" + randomId.String()
	bsConfig    conf.Bootstrap
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	// new zk client
	addr := bsConfig.Registrar.Zookeeper.Addr
	seconds := bsConfig.GetRegistrar().GetZookeeper().GetSessionTimeout().GetSeconds()
	nameSpace := bsConfig.GetRegistrar().GetZookeeper().GetNamespace()

	conn, _, err := zk.Connect(addr, time.Duration(seconds)*time.Second)
	if err != nil {
		panic(err)
	}

	zookeeper.WithRootPath("/" + nameSpace)
	// new reg with zk client
	reg := zookeeper.New(conn)

	return kratos.New(
		kratos.ID(instanceId),
		kratos.Name(bsConfig.GetServer().GetServiceName()),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(reg),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", instanceId,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bsConfig); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bsConfig.Server, bsConfig.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
