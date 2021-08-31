package oogway

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Start starts the Oogway server for given directory.
// The second argument is an optional template.FuncMap that will be merged into Oogway's funcmap.
func Start(dir string, funcmap template.FuncMap) error {
	ctx, cancel := context.WithCancel(context.Background())

	if err := watchConfig(ctx, dir); err != nil {
		cancel()
		return err
	}

	if err := loadPartials(dir); err != nil {
		cancel()
		return err
	}

	router := setupRouter(dir)

	if err := startServer(router, cancel); err != nil {
		return err
	}

	return nil
}

func setupRouter(dir string) *mux.Router {
	router := mux.NewRouter()
	serveAssets(router, dir)
	return router
}

func startServer(handler http.Handler, cancel context.CancelFunc) error {
	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		WriteTimeout: time.Second * time.Duration(cfg.HTTPWriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.HTTPReadTimeout),
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Shutting down server...")
		cancel()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.ShutdownTimeout))

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down server gracefully: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil

	/*if c.Server.HTTP.TLS {
		if err := server.ListenAndServeTLS(c.Server.HTTP.TLSCert, c.Server.HTTP.TLSKey); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	}*/
}
