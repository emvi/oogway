package oogway

import (
	"github.com/pirsch-analytics/pirsch-go-sdk"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	pirschClient *pirsch.Client
)

func initPirsch() {
	if cfg.Pirsch.ClientSecret != "" {
		pirschClient = pirsch.NewClient(cfg.Pirsch.ClientID, cfg.Pirsch.ClientSecret, "", nil)
		loadIPHeader()
		loadSubnets()
	}
}

func pageView(r *http.Request, path string) {
	if pirschClient != nil {
		url := r.URL

		if path != "" {
			url.Path = path
		}

		if err := pirschClient.HitWithOptions(r, &pirsch.HitOptions{
			IP:  getIP(r),
			URL: url.String(),
		}); err != nil {
			log.Printf("Error sending page view to Pirsch: %s", err)
		}
	}
}

func loadIPHeader() {
	for _, header := range cfg.Pirsch.Header {
		found := false

		for _, parser := range allIPHeader {
			if strings.ToLower(header) == strings.ToLower(parser.Header) {
				ipHeader = append(ipHeader, parser)
				found = true
				break
			}
		}

		if !found {
			log.Fatalf("Header invalid: %s", header)
		}
	}
}

func loadSubnets() {
	for _, subnet := range cfg.Pirsch.Subnets {
		_, n, err := net.ParseCIDR(subnet)

		if err != nil {
			log.Fatalf("Error parsing subnet '%s': %s", subnet, err)
		}

		allowedSubnets = append(allowedSubnets, *n)
	}
}
