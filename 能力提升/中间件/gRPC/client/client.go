package main

import (
	"context"
	"fmt"
	"log"
	pb "rpcDemo/proto"

	"google.golang.org/grpc"
)

func main() {
	// 连接服务器
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 实例化一个客户端
	client := pb.NewPersonServiceClient(conn)

	body := pb.Person{Name: "panyu"}

	// 调用接口
	resp, err := client.GetPersonInfo(context.Background(), &body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
