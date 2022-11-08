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

const (
	headline = `
 ________  ________  ________  ___       __   ________      ___    ___ 
|\   __  \|\   __  \|\   ____\|\  \     |\  \|\   __  \    |\  \  /  /|
\ \  \|\  \ \  \|\  \ \  \___|\ \  \    \ \  \ \  \|\  \   \ \  \/  / /
 \ \  \\\  \ \  \\\  \ \  \  __\ \  \  __\ \  \ \   __  \   \ \    / / 
  \ \  \\\  \ \  \\\  \ \  \|\  \ \  \|\__\_\  \ \  \ \  \   \/  /  /  
   \ \_______\ \_______\ \_______\ \____________\ \__\ \__\__/  / /    
    \|_______|\|_______|\|_______|\|____________|\|__|\|__|\___/ /     
                                                          \|___|/
v1.2.0`
)

var (
	baseDir    string
	tplFuncMap template.FuncMap
)

// Start starts the Oogway server for given directory.
// The second argument is an optional template.FuncMap that will be merged into Oogway's funcmap.
func Start(dir string, funcMap template.FuncMap) error {
	log.Println(headline)
	log.Printf("Starting in directory: %s", dir)
	baseDir = dir
	ctx, cancel := context.WithCancel(context.Background())

	if err := watchConfig(ctx, dir); err != nil {
		cancel()
		return err
	}

	initSass()

	if err := watchSass(ctx, dir); err != nil {
		cancel()
		return err
	}

	if err := watchJS(ctx, dir); err != nil {
		cancel()
		return err
	}

	tplFuncMap = mergeFuncMaps(funcMap)

	if err := watchPartials(ctx, dir, tplFuncMap); err != nil {
		cancel()
		return err
	}

	if err := watchContent(ctx, dir, tplFuncMap); err != nil {
		cancel()
		return err
	}

	initPirsch()
	router := setupRouter(dir)
	<-startServer(router, cancel)
	return nil
}

func setupRouter(dir string) *mux.Router {
	router := mux.NewRouter()
	serveSitemap(router)
	serveAssets(router, dir)
	servePage(router)
	return router
}

func startServer(handler http.Handler, cancel context.CancelFunc) chan struct{} {
	log.Println("Starting server...")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Handler:      handler,
		Addr:         addr,
		WriteTimeout: time.Second * time.Duration(cfg.Server.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.Server.ReadTimeout),
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Shutting down server...")
		cancel()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Server.ShutdownTimeout))

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down server gracefully: %s", err)
		}

		cancel()
	}()

	done := make(chan struct{})

	go func() {
		if cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
			if err := server.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile); err != nil {
				log.Fatalf("Error starting server: %s", err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %s", err)
			}
		}

		done <- struct{}{}
	}()

	log.Printf("Server started on http://%s!", addr)
	return done
}
