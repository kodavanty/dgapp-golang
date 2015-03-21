package stocks

import (
	"errors"
)

type Stock struct {
	Ticker   string  `json:"ticker"`
	Name     string  `json:"name"`
	Dividend float32 `json:"dividend"`
}

type StockManager struct {
	stocks []*Stock
}

func NewStockManager() *StockManager {
	return &StockManager{}
}

func NewStock(ticker string, name string, dividend float32) (s *Stock, err error) {
	return &Stock{ticker, name, dividend}, nil
}

func (sm *StockManager) Add(s *Stock) error {
	st := *s
	sm.stocks = append(sm.stocks, &st)

	return nil
}

func (sm *StockManager) All() (s []*Stock, err error) {
	return sm.stocks, nil
}

func (sm *StockManager) Find(ticker string) (s *Stock, err error) {
	for _, s := range sm.stocks {
		if s.Ticker == ticker {
			return s, nil
		}
	}

	return nil, errors.New("Could not find stock")
}
