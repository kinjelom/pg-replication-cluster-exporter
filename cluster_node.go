package main

import (
	"fmt"
)

type Node struct {
	host string
	db   *DataSource
}

type NodeState struct {
	host                   string
	err                    error
	isInRecovery           bool
	currentWalLsn          string
	currentWalLsnBytes     uint64
	lastWalReceiveLsn      string
	lastWalReceiveLsnBytes uint64
	lastWalReplayLsn       string
	lastWalReplayLsnBytes  uint64
}

func NewNode(db *DataSource, host string) *Node {
	return &Node{host: host, db: db}
}

func (n *Node) queryIsInRecovery() (bool, error) {
	var isInRecoveryStr string
	var err error
	if isInRecoveryStr, err = n.db.QueryStrWithEffort(n.host, "SELECT pg_is_in_recovery()"); err != nil {
		return false, fmt.Errorf("failed to query recovery mode: %v", err)
	}
	return isInRecoveryStr == "t" || isInRecoveryStr == "true", nil
}

func (n *Node) queryForState() *NodeState {
	var state = &NodeState{}
	state.host = n.host
	state.isInRecovery, state.err = n.queryIsInRecovery()
	if state.err == nil {
		// https://www.postgresql.org/docs/current/functions-admin.html
		if state.isInRecovery {
			// SLAVE
			if state.lastWalReceiveLsn, state.err = n.db.QueryStrWithEffort(n.host, "SELECT COALESCE(pg_last_wal_receive_lsn(),'0/0')"); state.err == nil {
				state.lastWalReceiveLsnBytes, state.err = parsePgLsn(state.lastWalReceiveLsn)
			} else {
				state.err = fmt.Errorf("failed to query last received wal location: %v", state.err)
			}
			if state.lastWalReplayLsn, state.err = n.db.QueryStrWithEffort(n.host, "SELECT COALESCE(pg_last_wal_replay_lsn(),'0/0')"); state.err == nil {
				state.lastWalReplayLsnBytes, state.err = parsePgLsn(state.lastWalReplayLsn)
			} else {
				state.err = fmt.Errorf("failed to query last replayed wal location: %v", state.err)
			}
		} else {
			// MASTER
			if state.currentWalLsn, state.err = n.db.QueryStrWithEffort(n.host, "SELECT COALESCE(pg_current_wal_lsn(),'0/0')"); state.err == nil {
				state.currentWalLsnBytes, state.err = parsePgLsn(state.currentWalLsn)
			} else {
				state.err = fmt.Errorf("failed to query current wal location: %v", state.err)
			}
		}
	}
	return state
}
