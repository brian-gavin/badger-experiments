package dump

import (
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

func Run(leaderRoot, followerRoot string) error {
	leader, err := badger.Open(badger.DefaultOptions(leaderRoot).WithReadOnly(true))
	if err != nil {
		return err
	}
	defer leader.Close()
	follower, err := badger.Open(badger.DefaultOptions(followerRoot).WithReadOnly(true))
	if err != nil {
		return err
	}
	defer follower.Close()
	fmt.Printf("LEADER at %s:\n", leaderRoot)
	if err := dump(leader); err != nil {
		return err
	}
	fmt.Printf("FOLLOWER at %s:\n", followerRoot)
	return dump(follower)
}
