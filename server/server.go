package main

import (
	"cangku/alliance"
	"cangku/db"
	pb "cangku/protobuf"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Cmd(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println(req)
	cmd := req.GetCmd()
	args := req.GetArgs()

	resp := &pb.Response{}

	switch cmd {
	case "whichAlliance":
	case "createAlliance":
		if len(args) != 2 {
			return resp, nil
		}
		am := &alliance.AllianceManager{}
		as := am.Create(args[0], args[1])
		if as == nil {
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)
	case "allianceList":
	case "joinAlliance":
	case "dismissAlliance":

	case "increaseCapacity":
		if len(args) != 1 {
			return resp, nil
		}
		as := &alliance.Alliance{}
		if err := as.AllianceStorageIncreaseCapacity(args[0], 10); err != nil {
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)

	case "storeItem":
		if len(args) != 4 {
			return resp, nil
		}
		allianceid := args[0]
		index, _ := strconv.Atoi(args[1])
		itemtype, _ := strconv.Atoi(args[2])
		count, _ := strconv.Atoi(args[3])
		as := &alliance.Alliance{}
		if err := as.AllianceStorageAddItem(allianceid, index, int32(itemtype), int32(count)); err != nil {
			fmt.Println(err)
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)
	case "destroyItem":

	case "clearup":
		if len(args) != 1 {
			return resp, nil
		}
		allianceid := args[0]
		as := &alliance.Alliance{}
		if err := as.AllianceStorageClearup(allianceid); err != nil {
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)
	}

	return resp, nil
}

func main() {
	db.InitRedisConnPool()

	lis, _ := net.Listen("tcp", ":12345")
	s := grpc.NewServer()

	pb.RegisterAllianceStorageServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("err: %+v", err)
	}
}
