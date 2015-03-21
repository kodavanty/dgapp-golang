package db

import (
	"github.com/gocql/gocql"
	"github.com/kodavanty/dgapp/stocks"
)

func InitDb(server string) (session *gocql.Session) {
	cluster := gocql.NewCluster(server)
	if cluster == nil {
		return nil
	}

	cluster.Keyspace = "stocks"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		return nil
	}

	return session
}

func AddDb(session *gocql.Session, stock *stocks.Stock) error {
	return session.Query(`INSERT INTO dg_stocks (ticker, name, dividend) VALUES (?, ?, ?)`,
		stock.Ticker, stock.Name, stock.Dividend).Exec()
}

func FindAllDb(session *gocql.Session) (s []stocks.Stock) {
	var ss [200]stocks.Stock

	iter := session.Query(`SELECT ticker, name, dividend FROM dg_stocks LIMIT 200`).Iter()
	i := 0
	for iter.Scan(&ss[i].Ticker, &ss[i].Name, &ss[i].Dividend) {
		if ss[i].Ticker == "" {
			break
		}
		i++
	}

	return ss[:i]
}

func FindDb(session *gocql.Session, ticker string) (e error, s *stocks.Stock) {
	var stock stocks.Stock

	err := session.Query(`SELECT ticker, name, dividend FROM dg_stocks WHERE ticker = ?`, ticker).Consistency(gocql.One).Scan(&stock.Ticker, &stock.Name, &stock.Dividend)

	return err, &stock
}
