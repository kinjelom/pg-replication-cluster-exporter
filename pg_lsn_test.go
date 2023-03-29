package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPgLsn_convertPgLsn(t *testing.T) {
	// pg_lsn supports values between 0/0 and FFFFFFFF/FFFFFFFF.

	// SELECT '0/0'::pg_lsn - '0/0'; 0
	val, _ := parsePgLsn("0/0")
	assert.Equal(t, uint64(0), val)

	// DEBUG> checking: master.currentWalLsn=`0/189B2E78`, slave.lastWalReceiveLsn=`0/90000A0`, slave.lastWalReplayLsn=`0/90000A0`
	// DEBUG> checking: slave.behind: `-261828056`, slave.delay: `0`

	// master.currentWalLsn=`0/189B2E78`
	// SELECT '0/189B2E78'::pg_lsn - '0/0'; -- 412823160
	val, _ = parsePgLsn("0/189B2E78")
	assert.Equal(t, uint64(412_823_160), val)

	// slave.lastWalReceiveLsn=`0/90000A0`, slave.lastWalReplayLsn=`0/90000A0`
	// SELECT '0/90000A0'::pg_lsn - '0/0'; -- 150995104
	val, _ = parsePgLsn("0/90000A0")
	assert.Equal(t, uint64(150_995_104), val)

	// SELECT '0/FFFFFFFF'::pg_lsn - '0/0'; -- 4294967295
	val, _ = parsePgLsn("0/FFFFFFFF")
	assert.Equal(t, uint64(4_294_967_295), val)

	// SELECT '7/A25801C8'::pg_lsn - '0/0'; -- 32788447688
	val, _ = parsePgLsn("7/A25801C8")
	assert.Equal(t, uint64(32_788_447_688), val)

	//  SELECT 'FFFFFFFF/0'::pg_lsn - '0/0'; -- 18446744069414584320
	val, _ = parsePgLsn("FFFFFFFF/0")
	assert.Equal(t, uint64(18_446_744_069_414_584_320), val)

	//  SELECT 'FFFFFFFF/FFFFFFFF'::pg_lsn - '0/0'; -- 18446744073709551615
	val, _ = parsePgLsn("FFFFFFFF/FFFFFFFF")
	assert.Equal(t, uint64(18_446_744_073_709_551_615), val)
}
