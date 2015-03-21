package db

import (
	"github.com/kodavanty/dgapp/stocks"
	"testing"
)

func TestDbInit(t *testing.T) {
	s := InitDb("localhost")
	if s == nil {
		t.Errorf("DB connection failed")
	}
	defer s.Close()
}

func TestDbAdd(t *testing.T) {
	s := InitDb("localhost")
	if s == nil {
		t.Errorf("DB connection failed")
	}
	defer s.Close()

	err := AddDb(s, &stocks.Stock{"ARUN", "ArubaNetworks", 0.0})
	if err != nil {
		t.Errorf("DB ADD failed %s", err)
	}

	err = AddDb(s, &stocks.Stock{"CSCO", "Cisco Systems", 1.5})
	if err != nil {
		t.Errorf("DB ADD failed %s", err)
	}

	err = AddDb(s, &stocks.Stock{"PG", "Proctor & Gamble", 3.8})
	if err != nil {
		t.Errorf("DB ADD failed %s", err)
	}

	err = AddDb(s, &stocks.Stock{"AAPL", "Applie Inc", 2.1})
	if err != nil {
		t.Errorf("DB ADD failed %s", err)
	}
}

func TestDbFind(t *testing.T) {
	s := InitDb("localhost")
	if s == nil {
		t.Errorf("DB connection failed")
	}
	defer s.Close()

	err, stock := FindDb(s, "ARUN")
	if err != nil {
		t.Errorf("DB FIND failed %s", err)
	}
	t.Logf("Stock %v", stock)

	err, stock = FindDb(s, "ABCD")
	if err == nil {
		t.Errorf("DB FIND failed %s", err)
	}
	t.Logf("Stock %v", stock)
}

func TestDbFindAll(t *testing.T) {
	s := InitDb("localhost")
	if s == nil {
		t.Errorf("DB connection failed")
	}
	defer s.Close()

	ss := FindAllDb(s)

	for _, st := range ss {
		t.Logf("Stock %v", st)
	}

	if len(ss) != 4 {
		t.Errorf("Count incorrect %d", len(ss))
	}
}
