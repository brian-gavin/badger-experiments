package main

import (
	"badger-experiments/backup"
	"context"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	"github.com/gogo/protobuf/proto"
)

func dump(db *badger.DB) error {
	s := db.NewStream()
	s.Send = func(kvs *pb.KVList) error {
		for _, kv := range kvs.GetKv() {
			proto.MarshalText(os.Stdout, kv)
		}
		return nil
	}
	return s.Orchestrate(context.Background())
}

func main() {
	leader, err := badger.Open(badger.DefaultOptions(backup.LeaderRoot).WithReadOnly(true))
	if err != nil {
		panic(err)
	}
	defer leader.Close()
	follower, err := badger.Open(badger.DefaultOptions(backup.FollowerRoot).WithReadOnly(true))
	if err != nil {
		panic(err)
	}
	defer follower.Close()
	fmt.Println("LEADER:")
	dump(leader)
	fmt.Println("FOLLOWER:")
	dump(follower)
}
