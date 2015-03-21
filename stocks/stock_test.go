package stocks

import (
	"testing"
)

func TestStockAdd(t *testing.T) {
	s, err := NewStock("ARUN", "ArubaNetworks", 2.0)
	sm := NewStockManager()
	err = sm.Add(s)
	if err != nil {
		t.Errorf("Stock addition failed")
	}
}

func TestStockFind(t *testing.T) {
	s, err := NewStock("CSCO", "Cisco Systems", 2.0)
	sm := NewStockManager()
	err = sm.Add(s)
	s, err = NewStock("ARUN", "ArubaNetworks", 2.0)
	err = sm.Add(s)

	_, err = sm.Find("NTAP")
	if err == nil {
		t.Errorf("Stock Find failed")
	}

	_, err = sm.Find("CSCO")
	if err != nil {
		t.Errorf("Stock Find failed")
	}

	ss, err1 := sm.All()
	if err1 != nil {
		t.Errorf("Stock Find failed")
	}
	if len(ss) != 2 {
		t.Errorf("Length failed")
	}
}

func TestStockCCC(t *testing.T) {
	ss := ParseCCCFile("/home/aruba/gocode/src/github.com/kodavanty/dgapp/USDIV.xlsx")

	t.Logf("Length: %d\n", len(ss))

	for _, s := range ss {
		t.Logf("Stock %s:%s:%f\n", s.Ticker, s.Name, s.Dividend)
	}

	if len(ss) != 105 {
		t.Errorf("Length failed")
	}
}
