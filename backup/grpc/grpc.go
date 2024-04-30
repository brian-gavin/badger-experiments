package grpc

//go:generate protoc --go_out=. --go-grpc_out=. ./backup.proto

import (
	"badger-experiments/backup/grpc/pb"
	"cmp"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	badgerPB "github.com/dgraph-io/badger/pb"
	protoV1 "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const (
	LeaderRoot   = "/tmp/backup/leader-grpc"
	FollowerRoot = "/tmp/backup/follower-grpc"
)

func Run() (err error) {
	if err = os.MkdirAll(LeaderRoot, 0o755); err != nil {
		return
	}
	if err = os.MkdirAll(FollowerRoot, 0o755); err != nil {
		return
	}
	stops := []chan bool{make(chan bool), make(chan bool)}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		for _, s := range stops {
			s <- true
		}
	}()
	go func() {
		err := leader(stops[0])
		if err != nil {
			panic(err)
		}
	}()
	err = follower(stops[1])
	return
}

// leaderHTTP starts an http server that accepts /set?key=<K>&value=<V>
func leaderHTTP() <-chan *http.Request {
	reqs := make(chan *http.Request, 100)
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.MarshalIndent(map[string]string{
			"inProgress": strconv.Itoa(len(reqs)),
		}, "", "  ")
		reqs <- r
		w.Write(resp)
	})
	go http.ListenAndServe(":8888", nil)
	return reqs
}

func leaderShipLog(c pb.BackupServiceClient, leader *badger.DB, lastBackup uint64) (uint64, error) {
	grpcStream, err := c.Backup(context.Background())
	if err != nil {
		return 0, err
	}
	dbStream := leader.NewStream()
	dbStream.LogPrefix = "LeaderShipLog"
	dbStream.ChooseKey = func(item *badger.Item) bool {
		return item.Version() >= lastBackup
	}
	var maxVersion uint64
	dbStream.Send = func(k *badgerPB.KVList) error {
		maxVersion = slices.MaxFunc(k.Kv, func(a, b *badgerPB.KV) int {
			return cmp.Compare(a.Version, b.Version)
		}).Version
		kvList, err := proto.Marshal(protoV1.MessageV2(k))
		if err != nil {
			return err
		}
		return grpcStream.Send(&pb.BackupRequest{EncodedKVs: kvList})
	}
	if err := dbStream.Orchestrate(context.Background()); err != nil {
		return 0, err
	}
	_, err = grpcStream.CloseAndRecv()
	return maxVersion, err
}

func leaderUpdate(leader *badger.DB, k, v string) {
	err := leader.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(k), []byte(v))
	})
	if err != nil {
		log.Printf("leader Update: %s", err)
	}
}

func leader(stop <-chan bool) (err error) {
	leader, err := badger.Open(badger.DefaultOptions(LeaderRoot))
	if err != nil {
		return
	}
	defer leader.Close()
	reqs := leaderHTTP()
	var lastBackup uint64
	tick := time.NewTicker(15 * time.Second)
	cc, err := grpc.NewClient("localhost:8484", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	c := pb.NewBackupServiceClient(cc)
	for {
		select {
		case <-tick.C:
			lastBackup, err = leaderShipLog(c, leader, lastBackup)
		case r := <-reqs:
			_ = r.Body.Close()
			leaderUpdate(leader, r.FormValue("key"), r.FormValue("value"))
		case <-stop:
			return
		}
	}
}

type followerServer struct {
	pb.UnimplementedBackupServiceServer
	db *badger.DB
}

func (s *followerServer) Backup(stream pb.BackupService_BackupServer) error {
	log.Println("follower: loading backup...")
	dbStreamWriter := s.db.NewStreamWriter()
	if err := dbStreamWriter.Prepare(); err != nil {
		return err
	}
	for {
		req, err := stream.Recv()
		log.Printf("follower: Recv error: %s", err)
		if err == io.EOF {
			log.Println("follower: done")
			err := dbStreamWriter.Flush()
			if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.BackupResponse{})
		}
		kvList := &badgerPB.KVList{}
		err = proto.Unmarshal(req.EncodedKVs, protoV1.MessageV2(kvList))
		if err != nil {
			log.Printf("follow load err: %s", err)
			return err
		}
		err = dbStreamWriter.Write(kvList)
		if err != nil {
			log.Printf("follow load stream.Writer err: %s", err)
			return err
		}
	}
}

func follower(stop <-chan bool) error {
	follower, err := badger.Open(badger.DefaultOptions(FollowerRoot))
	if err != nil {
		return err
	}
	defer follower.Close()
	server := &followerServer{db: follower}
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	defer grpcServer.GracefulStop()
	pb.RegisterBackupServiceServer(grpcServer, server)
	l, err := net.Listen("tcp", "localhost:8484")
	if err != nil {
		return err
	}
	go grpcServer.Serve(l)
	<-stop
	return nil
}
