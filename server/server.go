package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/colinsz/alliance/alliance"
	"github.com/colinsz/alliance/db"
	pb "github.com/colinsz/alliance/protobuf"

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
		if len(args) != 1 {
			resp.Result = "input like: whichAlliance playername"
			return resp, nil
		}
		am := &alliance.AllianceManager{}
		allinanceid, err := am.QueryByPlayer(args[0])
		if err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		resp.Result = string(allinanceid)
	case "createAlliance":
		if len(args) != 2 {
			resp.Result = "input like: whichAlliance alliancename playername"
			return resp, nil
		}
		am := &alliance.AllianceManager{}
		as, err := am.Create(args[0], args[1])
		if err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)
	case "allianceList":
		if len(args) != 0 {
			resp.Result = "input like: allianceList"
			return resp, nil
		}
		am := &alliance.AllianceManager{}
		alliancelist := am.List()

		b, _ := json.Marshal(alliancelist)
		resp.Result = "list: " + string(b)
	case "joinAlliance":
		if len(args) != 2 {
			resp.Result = "input like: joinAlliance alliancename playername"
			return resp, nil
		}
		allianceid := args[0]
		playername := args[1]
		am := &alliance.AllianceManager{}
		am.Join(playername, allianceid)
		resp.Result = "success"

	case "dismissAlliance":
		if len(args) != 2 {
			resp.Result = "input like: dismissAlliance alliancename playername"
			return resp, nil
		}
		allianceid := args[0]
		playername := args[1]
		am := &alliance.AllianceManager{}
		if err := am.Dismiss(allianceid, playername); err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		resp.Result = "success"

	case "increaseCapacity":
		if len(args) != 1 {
			resp.Result = "input like: increaseCapacity alliancename"
			return resp, nil
		}
		as := &alliance.Alliance{}
		err := as.AllianceStorageIncreaseCapacity(args[0], 10)
		if err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)

	case "storeItem":
		if len(args) != 4 {
			resp.Result = "input like: storeItem alliancename index itemtype count"
			return resp, nil
		}
		allianceid := args[0]
		index, _ := strconv.Atoi(args[1])
		itemtype, _ := strconv.Atoi(args[2])
		count, _ := strconv.Atoi(args[3])
		as := &alliance.Alliance{}
		err := as.AllianceStorageAddItem(allianceid, index, int32(itemtype), int32(count))
		if err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)
	case "destroyItem":
		if len(args) != 3 {
			resp.Result = "input like: destroyItem alliancename playername index"
			return resp, nil
		}
		allianceid := args[0]
		playername := args[1]
		index, _ := strconv.Atoi(args[2])
		as := &alliance.Alliance{}
		err := as.AllianceStorageDestroyItem(allianceid, playername, index)
		if err != nil {
			resp.Result = err.Error()
			return resp, nil
		}
		b, _ := json.Marshal(as)
		resp.Result = string(b)

	case "clearup":
		if len(args) != 1 {
			resp.Result = "input like: clearup alliancename"
			return resp, nil
		}
		allianceid := args[0]
		as := &alliance.Alliance{}
		err := as.AllianceStorageClearup(allianceid)
		if err != nil {
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
