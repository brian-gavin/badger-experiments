# badger-experiments

Scientific experiments[^1] with [https://dgraph.io/docs/badger/](Badger).

## Backup

This example found in `/backup` uses the `Backup` and `Load` APIs to do simple replication via
shipping the leader's DB over TCP to the follower on an interval. A more robust solution would do
many other things, and also use the underlying `Stream` API for more control over the send logic.

The backup will open two DB's at `/tmp/backup/leader` and `/tmp/backup/follower`. The leader will
http listen on port `8888` for updates to the DB using a simple api: `/set?K=<k>&V=<v>`. The
follower will tcp listen port `8282` for backups sent from the leader.

`CTRL-C` cleanly stops the program.

### To run the Backup experiment

```
# In one terminal
mkdir -p /tmp/backup/{leader,follower}
go run ./backup/example

# In another terminal
curl 'localhost:8888/set?key=K&value=V' # set's K=V in the Leader's DB.
```

`CTRL-C` the `example` process to stop the leader and follower, and use the `dump` program to
inspect the leader and follower's DBs.

```
go run ./backup/dump
```

[^1]: https://www.youtube.com/watch?v=BSUMBBFjxrY
