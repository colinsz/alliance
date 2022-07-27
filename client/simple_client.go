package main

import (
	"cangku/protobuf"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
)

const addr = "127.0.0.1:12345"

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) <= 1 {
		return
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial: %+v", err)
		return
	}
	defer conn.Close()

	cli := protobuf.NewAllianceStorageClient(conn)
	req := protobuf.Request{
		Cmd:  args[0],
		Args: args[1:],
	}
	resp, err := cli.Cmd(context.Background(), &req)
	if err != nil {
		log.Fatalf("cmd: %+v", err)
		return
	}
	log.Println(resp)
}
