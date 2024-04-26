package backup

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	LeaderRoot   = "/tmp/backup/leader"
	FollowerRoot = "/tmp/backup/follower"
)

func Run() (err error) {
	stops := []chan bool{make(chan bool), make(chan bool)}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		for _, s := range stops {
			s <- true
		}
	}()
	listening := make(chan bool)
	go follower(stops[1], listening)
	<-listening
	err = leader(stops[0])
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

func leaderShipLog(leader *badger.DB, lastBackup uint64) (uint64, error) {
	fconn, err := net.Dial("tcp", "127.0.0.1:8484")
	if err != nil {
		return 0, err
	}
	defer fconn.Close()
	lastBackup, err = leader.Backup(fconn, lastBackup)
	if err != nil {
		log.Printf("leader backup: %s", err)
	}
	return lastBackup, err
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
	for {
		select {
		case <-tick.C:
			lastBackup, err = leaderShipLog(leader, lastBackup)
		case r := <-reqs:
			_ = r.Body.Close()
			leaderUpdate(leader, r.FormValue("key"), r.FormValue("value"))
		case <-stop:
			return
		}
	}
}

func followerGoAcceptLoop(l net.Listener) <-chan net.Conn {
	conns := make(chan net.Conn)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				panic(err)
			}
			conns <- c
		}
	}()
	return conns
}

func follower(stop <-chan bool, listening chan<- bool) {
	l, err := net.Listen("tcp", "127.0.0.1:8484")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	listening <- true
	conns := followerGoAcceptLoop(l)
	follower, err := badger.Open(badger.DefaultOptions(FollowerRoot))
	if err != nil {
		panic(err)
	}
	defer follower.Close()
	for {
		var c net.Conn
		select {
		case <-stop:
			return
		case c = <-conns:
		}
		log.Println("follower: loading backup...")
		err = follower.Load(c, 1024)
		if err != nil {
			log.Printf("follow load err: %s", err)
		}
		log.Println("follower: done")
		c.Close()
	}
}
