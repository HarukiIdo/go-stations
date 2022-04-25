package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	// .envに環境変数を書いている？
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// .envに環境変数を書いている？
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}

	// 関数の遅延実行
	defer todoDB.Close()

	// set http handlers
	mux := router.NewRouter(todoDB)

	// start the http server
	error := http.ListenAndServe(port, mux)

	return error
}
