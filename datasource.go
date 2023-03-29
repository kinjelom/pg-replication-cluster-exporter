package main

import (
	"database/sql"
	"fmt"
	"time"
)

type DataSource struct {
	measurer   *Measurer
	driverName string
	port       string
	dbname     string
	user       string
	password   string
	sslMode    string
	connection map[string]*sql.DB
}

func NewDataSource(measurer *Measurer, port, user, password string) *DataSource {
	return &DataSource{
		measurer:   measurer,
		driverName: "postgres",
		port:       port,
		dbname:     "postgres",
		user:       user,
		password:   password,
		sslMode:    "disable",
		connection: make(map[string]*sql.DB),
	}
}

func (db *DataSource) connect(host string, force bool) (*sql.DB, error) {
	var err error
	if db.connection[host] == nil || force {
		cs := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, db.port, db.dbname, db.user, db.password, db.sslMode)
		db.connection[host], err = sql.Open(db.driverName, cs)
		if err != nil {
			db.connection[host] = nil
			log.warn("Can't connect to %s, error: %v", host, err)
			return nil, err
		}
	}
	return db.connection[host], nil
}

func (db *DataSource) reconnect(host string) (*sql.DB, error) {
	db.measurer.incReconnects(host)
	var err error
	_, err = db.connect(host, true)
	if err != nil {
		db.connection[host] = nil
		log.warn("Can't connect %s, error: %v", host, err)
		return nil, err
	} else {
		err = db.connection[host].Ping()
		if err != nil {
			db.connection[host] = nil
			log.warn("Can't ping %s, error: %v", host, err)
			return nil, err
		}
	}
	return db.connection[host], nil
}

func (db *DataSource) Ping(host string) (int64, error) {
	start := time.Now()
	if db.connection[host] == nil {
		return -1, fmt.Errorf("host %s is disconnected", host)
	}
	err := db.connection[host].Ping()
	return time.Since(start).Milliseconds(), err
}

func (db *DataSource) QueryOneValWithEffort(host, q string) (string, error) {
	log.debug("query: `%s`", q)
	var err error
	_, err = db.connect(host, false)
	if err != nil {
		log.warn("Can't connect %s, error: %v", host, err)
		return "", err
	}
	var v string
	v, err = db.queryOneVal(host, q)
	if err != nil {
		_, err = db.reconnect(host)
		if err != nil {
			log.warn("Can't connect %s, error: %v", host, err)
			return "", err
		} else {
			v, err = db.queryOneVal(host, q)
		}
	}
	log.debug("query result: `%s`", v)
	return v, err
}

func (db *DataSource) queryOneVal(host, q string) (string, error) {
	start := time.Now()
	row := db.connection[host].QueryRow(q)
	var v string
	if err := row.Scan(&v); err != nil {
		db.measurer.updateQueryStats(host, q, time.Since(start).Milliseconds(), false)
		return "", err
	}
	db.measurer.updateQueryStats(host, q, time.Since(start).Milliseconds(), true)
	return v, nil
}
