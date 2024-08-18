package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		ErrorLog:     log.New(app.logger, "", 0),
	}
	// graceful shutdown
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		c := <-quit
		app.logger.PrintInfo("shuting down the server", map[string]string{
			"signal": c.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := ser.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.PrintInfo("completing background tasks",
			map[string]string{"addr": ser.Addr})
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("Lisetening on port: ", map[string]string{
		"addr": ser.Addr,
	})
	if err := ser.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
