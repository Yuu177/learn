package main

import (
	"context"
	"log"
	"net"

	pb "rpcDemo/proto"

	"google.golang.org/grpc"
)

type PersonService struct {
}

func (p *PersonService) GetPersonInfo(ctx context.Context, req *pb.Person) (*pb.Person, error) {
	if req.Name == "panyu" {
		return &pb.Person{
			Name:    "panyu",
			Age:     18,
			Address: []string{"china", "usa"},
		}, nil
	}
	return nil, nil
}

func main() {
	// 监听
	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}

	// 实例化 grpc
	s := grpc.NewServer()

	// gRPC 上注册服务
	p := PersonService{}
	pb.RegisterPersonServiceServer(s, &p) // 注册 PersonService struct 下所有的方法

	// 启动服务器
	s.Serve(listener)
}
