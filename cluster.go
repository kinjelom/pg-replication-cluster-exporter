package main

import (
	"fmt"
)

type Cluster struct {
	name       string
	nodes      map[string]*Node
	dataSource *DataSource
}

type SlaveLag struct {
	receiveLag uint64
	replayLag  uint64
}

func NewCluster(dataSource *DataSource, clusterName string, hosts []string) *Cluster {
	cluster := &Cluster{}
	cluster.name = clusterName
	cluster.dataSource = dataSource
	cluster.nodes = make(map[string]*Node)
	for _, host := range hosts {
		cluster.nodes[host] = NewNode(cluster.dataSource, host)
	}
	return cluster
}

func (cluster *Cluster) queryForState() (*NodeState, *map[string]*NodeState, error) {
	var master = &NodeState{}
	var slaves = make(map[string]*NodeState)
	mi := 0
	si := 0
	for host, node := range cluster.nodes {
		nodeState := node.queryForState()
		if nodeState.err == nil {
			if nodeState.isInRecovery {
				slaves[host] = nodeState
				si++
			} else {
				if mi == 0 {
					master = nodeState
					mi++
				} else {
					return nil, nil, fmt.Errorf("too many masters, konwn %s, pretending: %s", master.host, node.host)
				}
			}
		}
	}
	if mi == 0 || si == 0 {
		return master, &slaves, fmt.Errorf("this is not replication cluster, masters: %d, slaves: %d", mi, si)
	}
	return master, &slaves, nil
}

func (cluster *Cluster) calculateSlaveLag(master NodeState, slave NodeState) *SlaveLag {
	lag := &SlaveLag{receiveLag: 0, replayLag: 0}
	if master.currentWalLsnBytes > slave.lastWalReceiveLsnBytes {
		lag.receiveLag = master.currentWalLsnBytes - slave.lastWalReceiveLsnBytes
	} else {
		lag.receiveLag = 0
	}
	if slave.lastWalReceiveLsnBytes > slave.lastWalReplayLsnBytes {
		lag.replayLag = slave.lastWalReceiveLsnBytes - slave.lastWalReplayLsnBytes
	} else {
		lag.replayLag = 0
	}
	log.debug("calculate lag between master slave %s:\n"+
		"  master.currentWalLsn    = %d (%s)\n"+
		"  slave.lastWalReceiveLsn = %d (%s)\n"+
		"  slave.lastWalReplayLsn  = %d (%s)\n"+
		"  slave.receiveLag        = %d\n"+
		"  slave.replayLag         = %d",
		slave.host,
		master.currentWalLsnBytes, master.currentWalLsn,
		slave.lastWalReceiveLsnBytes, slave.lastWalReceiveLsn,
		slave.lastWalReplayLsnBytes, slave.lastWalReplayLsn,
		lag.receiveLag, lag.replayLag)
	return lag
}
