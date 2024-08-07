package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/jsonlog"
	_ "github.com/lib/pq"
)

// all configurations the application needs to run
type config struct {
	port int
	db   struct {
		dsn          string
		maxIdleTime  string
		maxIdleConns int
		maxOpenConns int
	}
}

type application struct {
	config config
	models data.Models
	logger *jsonlog.Logger
}

func main() {
	config := config{}

	// server port
	flag.IntVar(&config.port, "port", 4000, "server port")
	// database configuration
	flag.StringVar(&config.db.dsn, "db-dsn", "", "database dsn")
	flag.StringVar(&config.db.maxIdleTime,
		"db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.IntVar(&config.db.maxOpenConns,
		"db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&config.db.maxIdleConns,
		"db-max-idle-conns", 25, "PostgreSQL max idle connections")

	flag.Parse()
	logg := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logg.PrintInfo("database connection pool established", nil)
	// initiate application
	app := &application{
		config: config,
		models: data.NewModels(db),
		logger: logg,
	}
	if err := app.serve(); err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

func openDB(conf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(conf.db.maxIdleConns)
	db.SetMaxOpenConns(conf.db.maxOpenConns)
	maxDuration, err := time.ParseDuration(conf.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(maxDuration)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
