package oogway

import (
	"encoding/xml"
	"time"
)

const (
	sitemapLastModFormat = "2006-01-02"
	header               = `<?xml version="1.0" encoding="UTF-8"?>`
	xmlns                = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type sitemapURLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []sitemapURL
}

type sitemapURL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	Lastmod    string   `xml:"lastmod"`
	Changefreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

func generateSitemap(urls []sitemapURL) ([]byte, error) {
	now := time.Now().Format(sitemapLastModFormat)

	for i := range urls {
		if urls[i].Lastmod == "" {
			urls[i].Lastmod = now
		}
	}

	sitemap := sitemapURLSet{
		XMLNS: xmlns,
		URLs:  urls,
	}
	out, err := xml.Marshal(&sitemap)

	if err != nil {
		return nil, err
	}

	return []byte(header + string(out)), nil
}
