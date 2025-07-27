package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	_ "github.com/piitschy/qrcode-me/migrations"
	"github.com/piitschy/qrcode-me/proxies"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {

	sttl := os.Getenv("CACHE_TTL")
	if sttl == "" {
		sttl = "60s" // default to 60 seconds if not set
	}

	ttl, err := time.ParseDuration(sttl)
	if err != nil {
		panic("Invalid CACHE_TTL value: " + sttl)
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{slug}", func(e *core.RequestEvent) error {
			slug := e.Request.PathValue("slug")
			if slug == "" {
				return e.Error(http.StatusNotFound, "Page not found", nil)
			}

			url, found := cache.Get(slug)
			if found {
				e.Redirect(http.StatusTemporaryRedirect, url)
				return nil
			}

			link := &proxies.Link{}
			err := app.RecordQuery("links").
				Where(dbx.HashExp{"status": "public"}).
				AndWhere(dbx.NewExp("slug = {:slug}", dbx.Params{"slug": slug})).
				One(link)
			if err != nil {
				return e.Error(http.StatusNotFound, "Page not found", err)
			}
			// Increment the click count
			link.IncrementCount()
			if err := app.Save(*link); err != nil {
				app.Logger().Error("Failed to update link count", "error", err)
			}
			cache.SetWithTTL(slug, link.Url(), 1, ttl)
			e.Redirect(http.StatusTemporaryRedirect, link.Url())
			return nil
		})

		return se.Next()
	})

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
