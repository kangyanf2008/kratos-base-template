package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/zookeeper/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-zookeeper/zk"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/anypb"
	providerV1 "kratos-base-template/api/provider/v1"
	"os"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
	"kratos-base-template/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = ""
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _    = os.Hostname()
	bsConfig conf.Bootstrap
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config-client.yaml")
}

func main() {
	flag.Parse()
	/*	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)*/
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

	addr := bsConfig.Registrar.Zookeeper.Addr
	seconds := bsConfig.GetRegistrar().GetZookeeper().GetSessionTimeout().GetSeconds()
	nameSpace := bsConfig.GetRegistrar().GetZookeeper().GetNamespace()

	conn, _, err := zk.Connect(addr, time.Duration(seconds)*time.Second)
	if err != nil {
		panic(err)
	}

	zookeeper.WithRootPath("/" + nameSpace)
	// new reg with zk client
	discovery := zookeeper.New(conn)
	endpoint := "discovery:///" + bsConfig.GetServer().GetServiceName()
	// 创建 P2C 负载均衡算法 Selector，并将路由 Filter 注入
	selector.SetGlobalSelector(wrr.NewBuilder())
	//selector.SetGlobalSelector(p2c.NewBuilder())
	clientConn, err2 := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(discovery),
	)

	if err2 != nil {
		panic(err2)
	}
	defer clientConn.Close()
	client := providerV1.NewProviderClient(clientConn)

	var i = 0
	begin := time.Now().Unix()
	for true {
		i++
		body, _ := anypb.New(&wrappers.StringValue{Value: "小明#########################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################" +
			"#####################################################################"})

		request := &providerV1.Request{
			ReqId: strconv.Itoa(i),
			Event: providerV1.EVENT_CODE_ADD_USER,
			Body:  body,
		}

		_, err3 := client.BaserService(context.Background(), request)
		if err3 != nil {
			fmt.Printf("erro=%v \n", err3)
		}
		if i%100000 == 0 {
			endTime := time.Now().Unix()
			fmt.Println(strconv.FormatInt(endTime-begin, 10) + "s")
			begin = endTime
		}
		//fmt.Printf("response=%v, erro=%v \n", response, err3)
		//time.Sleep(time.Second * 1)
	}

}
