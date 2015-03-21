package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kodavanty/dgapp/db"
	"github.com/kodavanty/dgapp/stocks"
)

type DGAppConfig struct {
	LogPath  string
	DbServer string
	CccFile  string
}

type DGAppInfo struct {
	config DGAppConfig
	sm     *stocks.StockManager
	db     *gocql.Session
}

var (
	AppInfo DGAppInfo
	MyLog   *log.Logger
)

func panic_on_err(err error) {
	if err != nil {
		panic(err)
	}
}

/* Read parameters from the application config file */
func readConfig() {
	configFile := os.Getenv("DGAPP_CONFIG")
	if configFile == "" {
		fmt.Println("Config file not found")
		return
	}
	fmt.Printf("Reading from config file %s\n", configFile)

	bytes, err := ioutil.ReadFile(configFile)
	panic_on_err(err)

	err = yaml.Unmarshal(bytes, &AppInfo.config)
	panic_on_err(err)

	fmt.Printf("Config %v\n", AppInfo.config)
}

/* Setup the logging to file */
func setupLogger() {
	f, err := os.OpenFile(AppInfo.config.LogPath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		os.FileMode(0666))
	defer f.Close()
	panic_on_err(err)

	MyLog = log.New(f, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getTicker(r *http.Request) (string, error) {
	txt := mux.Vars(r)["ticker"]

	if txt == "" {
		return "", errors.New("Ticker argument not found")
	}

	return txt, nil
}

func getAllStockHandler(w http.ResponseWriter, r *http.Request) error {
	Stocks, err := AppInfo.sm.All()
	panic_on_err(err)

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(Stocks)
}

func getStockHandler(w http.ResponseWriter, r *http.Request) error {
	ticker, e := getTicker(r)
	if e != nil {
		return errors.New("Could not get ticker")
	}

	s, e1 := AppInfo.sm.Find(ticker)
	if e1 != nil {
		MyLog.Printf("Finding stock %s from DB\n", ticker)
		e1, s = db.FindDb(AppInfo.db, ticker)
		if e1 != nil {
			return errors.New("Could not get stock")
		}
		AppInfo.sm.Add(s)
		MyLog.Printf("Found stock %s from DB\n", ticker)
	}
	MyLog.Println(s)

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(*s)
}

func apiHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Path %s!\n", r.URL.Path)

	return nil
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

/* Setup webserver and routing */
func setupWebServer() {
	/* Setup routing for URI's */
	r := mux.NewRouter()
	r.HandleFunc("/api/", errorHandler(getAllStockHandler)).Methods("GET")
	r.HandleFunc("/api/{ticker}", errorHandler(getStockHandler)).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	r.PathPrefix("/js/{rest}").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/css/{rest}").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/html/{rest}").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("static"))))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func InitDbTemp() {
	s, err := stocks.NewStock("ARUN", "ArubaNetworks", 0.0)
	panic_on_err(err)
	err = AppInfo.sm.Add(s)
	panic_on_err(err)

	s, err = stocks.NewStock("CSCO", "Cisco Systems", 3.4)
	panic_on_err(err)
	err = AppInfo.sm.Add(s)
	panic_on_err(err)

	s, err = stocks.NewStock("AAPL", "Apple Inc", 1.5)
	panic_on_err(err)
	err = AppInfo.sm.Add(s)
	panic_on_err(err)
}

func Init() {
	readConfig()
	setupLogger()

	AppInfo.db = db.InitDb(AppInfo.config.DbServer)
	if AppInfo.db == nil {
		panic(errors.New("DB NULL"))
	}

	//Parse the CCC list
	ss := stocks.ParseCCCFile(AppInfo.config.CccFile)
	AppInfo.sm = stocks.NewStockManager()

	//Add it to the DB
	var err error
	for _, s := range ss {
		err = db.AddDb(AppInfo.db, &s)
		panic_on_err(err)
	}

	//Cache everything
	for _, s := range db.FindAllDb(AppInfo.db) {
		err = AppInfo.sm.Add(&s)
		panic_on_err(err)
	}

	setupWebServer()
}
