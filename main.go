package main

import (
	"context"
	"executor/pb"
	"executor/utils"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type judgeServer struct {}
func (j *judgeServer) Judge(ctx context.Context, jc *pb.JudgeConfig) (*pb.JudgeResult, error) {
	fmt.Println("ok this test success")
	return &pb.JudgeResult{
		JudgeId: jc.JudgeId,
	}, nil
}

func main() {
	var (
		id = flag.String("id", utils.UUID(10), "实例id")
		address = flag.String("address", "127.0.0.1", "服务地址")
		port = flag.Int("port", 12100, "服务端口")
	)
	flag.Parse()
	if err := consulRegister(*id, *address, *port); err != nil {
		fmt.Errorf("consul registered failed, err: %v\n", err)
		return
	}

	if err := grpcServe(*address, *port); err != nil {
		fmt.Errorf("grpc serve failed, err: %v\n", err)
		return
	}
}

func consulRegister(id string, address string, port int) error {
	// 将grpc服务注册到consul上
	// 1、初始化consul配置
	consulConfig := api.DefaultConfig()
	// 2、创建consul对象
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("api.NewClient err:", err)
		return err
	}
	// 3、告诉consul，即将注册的服务的配置信息
	reg := api.AgentServiceRegistration{
		ID:      id,
		Tags:    []string{"exec"},
		Name:    "ytoj-exec",
		Address: address,
		Port:    port,
		Check: &api.AgentServiceCheck{
			CheckID:  id,
			TCP:      address + ":" + strconv.Itoa(port),
			Timeout:  "10s",
			Interval: "5s",
		},
	}
	// 4、注册grpc服务到consul
	err = consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		fmt.Println("register err:", err)
		return err
	}
	return nil
}

func grpcServe(address string, port int) error {
	// 1、初始化grpc对象
	grpcServer := grpc.NewServer()
	// 2、注册服务
	pb.RegisterJudgeServer(grpcServer, new(judgeServer))
	// 3、设置监听，指定ip port
	ln, err := net.Listen("tcp", address + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("listen err: ", err)
		return err
	}
	defer ln.Close()
	// 4、启动服务
	fmt.Println("start serve...")
	err = grpcServer.Serve(ln)
	if err != nil {
		fmt.Println("serve err: ", err)
		return err
	}
	return nil
}