package oogway

import (
	"github.com/pirsch-analytics/pirsch-go-sdk"
	"log"
	"net/http"
)

var (
	pirschClient *pirsch.Client
)

func initPirsch() {
	if cfg.Pirsch.ClientSecret != "" {
		pirschClient = pirsch.NewClient(cfg.Pirsch.ClientID, cfg.Pirsch.ClientSecret, "", nil)
	}
}

func pageView(r *http.Request, path string) {
	if pirschClient != nil {
		url := r.URL

		if path != "" {
			url.Path = path
		}

		if err := pirschClient.HitWithOptions(r, &pirsch.HitOptions{
			// TODO add support for IP header
			URL: url.String(),
		}); err != nil {
			log.Printf("Error sending page view to Pirsch: %s", err)
		}
	}
}
