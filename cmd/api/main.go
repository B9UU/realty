package main

import (
	"flag"
	"fmt"
	"log"
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
}

func main() {
	config := config{}
	flag.IntVar(&config.port, "port", 4000, "server port")
	flag.Parse()
	fmt.Println(config.port)
	app := &application{
		config: config,
	}
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

}
