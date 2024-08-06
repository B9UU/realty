package main

import (
	"fmt"
	"net/http"
	"time"
)

// server config and serve..
func (app *application) serve() error {
	ser := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		// TODO: add error log handler
	}
	app.logger.PrintInfo("Lisetening on port: ", map[string]string{
		"addr": ser.Addr,
	})
	if err := ser.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
