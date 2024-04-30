package main

import (
	"badger-experiments/backup/dump"
	"badger-experiments/backup/grpc"
	"badger-experiments/backup/tcp"
	"flag"
	"fmt"
	"os"
)

var (
	example = flag.String("example", "", "name of the backup example to run (tcp or grpc)")
	runDump = flag.Bool("dump", false, "dump the leader and follower DBs to terminal")
)

func main() {
	flag.Parse()
	var (
		runner                   func() error
		leaderRoot, followerRoot string
	)
	switch e := *example; e {
	case "tcp":
		runner = tcp.Run
		leaderRoot, followerRoot = tcp.LeaderRoot, tcp.FollowerRoot
	case "grpc":
		runner = grpc.Run
		leaderRoot, followerRoot = grpc.LeaderRoot, grpc.FollowerRoot
	case "":
		fmt.Fprintln(os.Stderr, "-example must be set")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s is an invalid -example\n", e)
		os.Exit(1)
	}
	var err error
	if *runDump {
		err = dump.Run(leaderRoot, followerRoot)
	} else {
		err = runner()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
