package main

import (
	"fmt"
	"strconv"
	"strings"
)

// https://pgpedia.info/p/pg_lsn.html
func parsePgLsn(s string) (uint64, error) {
	const supportInfo = "parser supports values between 0/0 and FFFFFFFF/FFFFFFFF"
	if s == "" {
		s = "0/0"
	}
	l := strings.SplitN(s, "/", 3)
	if len(l) != 2 {
		return 0, fmt.Errorf("parsing pg_lsn=`%s` failed (not two parts?), %s", s, supportInfo)
	}

	a, err := strconv.ParseUint(l[0], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing pg_lsn=`%s` failed (first part), %s", s, supportInfo)
	}

	b, err := strconv.ParseUint(l[1], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing pg_lsn=`%s` failed (second part), %s", s, supportInfo)
	}

	return a<<32 + b, nil
}
