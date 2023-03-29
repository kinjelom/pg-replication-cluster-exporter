# Postgres Replication Cluster Exporter

This repository houses `pgrc_exporter`, a small utility that will connect to all the nodes in a Postgres replication cluster and exports replication lag metrics of each node (Prometheus format).

Inspired by: https://github.com/Qarik-Group/pg-replication-tester - thanks :+1:

## Metrics

- **build_info**: Program build info
- **cluster_node_info**: Cluster node info
- **reconnects_count_total**: Cluster node reconnects total count
- **queries_count_total**: All queries total count
- **last_query_seconds**: Cluster node last query seconds
- **current_wal_lsn_bytes**: `SELECT pg_current_wal_lsn()`
- **last_wal_receive_lsn_bytes**: `SELECT pg_last_wal_receive_lsn()`
- **last_wal_replay_lsn_bytes**: `SELECT pg_last_wal_replay_lsn()`
- **receive_lag_bytes**: Cluster node receive bytes: `pg_current_wal_lsn() - pg_last_wal_receive_lsn()`
- **replay_lag_bytes**: Cluster node replay bytes: `pg_last_wal_receive_lsn() - pg_last_wal_reply_lsn()`

## Options

```
-A, --address, Address to listens on the TCP network. Default: :9188
-P, --path, Path under which to expose metrics. Default: /metrics
-C, --cluster-name, Cluster name. Default: cluster-hash(nodes)
-n, --node, Replication cluster nodes. May be specified more than once.
-p, --port, TCP port that Postgres listens on. Default: 6432 
-u, --user, User to connect as.
-s, --password, Password to connect with.
-i, --interval, Collecting metrics interval in seconds. Default: 15 
-V, --verbosity, Verbosity level (0 errors, 1 +warnings, 2 +infos, 3 +debugs). Default: 2 
-v, --version, Output version information, then exit.
-h, --help, Show this help, then exit.
```

## Building

```bash
GOOS=linux GOARCH=amd64 go build -o pgrc_exporter
```
