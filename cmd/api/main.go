package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/scraper"
	"github.com/b9uu/realty/internal/validator/mailer"
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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	models data.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
	mailer mailer.Mailer
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

	// mail config
	flag.StringVar(&config.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&config.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&config.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&config.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&config.smtp.sender,
		"smtp-sender", "Realty <noreply@ibrahimboussaa.com>", "SMTP sender")

	// rate limiter config
	flag.Float64Var(&config.limiter.rps,
		"limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()
	logg := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logg.PrintInfo("database connection pool established", nil)
	// initiate application
	fmt.Printf("%+v", config.smtp)
	app := &application{
		config: config,
		models: data.NewModels(db),
		logger: logg,
		mailer: mailer.New(
			config.smtp.host, config.smtp.sender,
			config.smtp.username, config.smtp.password,
			config.smtp.port,
		),
	}
	err = scraper.Scrape(app.models.Realty)
	if err != nil {
		panic(err)
	}

	var version = "1.0.0"
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
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
