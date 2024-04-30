# badger-experiments

Scientific experiments[^1] with [https://dgraph.io/docs/badger/](Badger).

## Backup

There are two Backup experiments: one that uses the `Backup` and `Load` APIs to do simple replication
over TCP. The second uses the `Stream` and `StreamWriter` APIs to do simple replication over gRPC.
They are called the TCP and gRPC examples respectively. 

The backup will open two DB's at `/tmp/backup-{tcp,grpc}/leader` and `/tmp/backup-{tcp,grpc}/follower`.
The leader will http listen on port `8888` for updates to the DB using a simple api: `/set?K=<k>&V=<v>`.
The follower will tcp listen port `8484` for backups sent from the leader.

### To run the Backup experiments

Experiments are run via `go run ./backup` with flags to control which backup is to be run.

`CTRL-C` cleanly stops the backup programs.

```
# To run the TCP experiment
go run ./backup -example tcp

# To run the gRPC experiment
go run ./backup -example grpc

# To update the leader DB
curl 'localhost:8888/set?key=K&value=V' # set's K=V in the Leader's DB.
```

`CTRL-C` the `backup` process to stop the leader and follower, and use the `dump` program to
inspect the leader and follower's DBs to ensure that the replication occurred correctly.
The `-example` flag must be used to ensure that the dump program reads from the right database.

```
go run ./backup/ -example grpc -dump
```

The example found in `/backup/tcp` uses the `Backup` and `Load` APIs to do simple replication via
shipping the leader's DB over TCP to the follower on an interval.

The `/backup/grpc` experiment uses the `Stream` and `StreamWriter` APIs to experiment with a more
accurate solution. It uses a gRPC client stream to send a large batch from a `Stream` (on the leader)
to a `StreamWriter` (on the follower), because this seems like the right thing to do at the time.

[^1]: https://www.youtube.com/watch?v=BSUMBBFjxrY
