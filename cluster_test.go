package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCluster_calculateLag(t *testing.T) {
	cwLsn, _ := parsePgLsn("0/189B2E78") // 412_823_160
	m := NodeState{
		host:               "testMaster",
		currentWalLsnBytes: cwLsn,
	}
	lwrcLsn, _ := parsePgLsn("0/90000A1") // 150_995_105
	lwrpLsn, _ := parsePgLsn("0/90000A0") // 150_995_104
	s := NodeState{
		host:                   "testSlave",
		lastWalReceiveLsnBytes: lwrcLsn,
		lastWalReplayLsnBytes:  lwrpLsn,
	}
	cluster := &Cluster{}
	lag := cluster.calculateSlaveLag(m, s)

	assert.Equal(t, uint64(261_828_055), lag.receiveLag)
	assert.Equal(t, uint64(1), lag.replayLag)
}
