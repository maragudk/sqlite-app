package http

import (
	"log"
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

// RedirectRegion using the fly-replay HTTP header if "region" is set in the URL query params.
func RedirectRegion(currentRegion string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if region := r.URL.Query().Get("region"); region != "" && region != currentRegion {
				w.Header().Set("fly-replay", "region="+region)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type primaryGetter interface {
	GetPrimary() (string, error)
}

// RedirectToPrimary if the request is POST and this is not the primary instance.
func RedirectToPrimary(db primaryGetter, log *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We don't need to do anything if this is not a POST request
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}

			// If region is forced in the URL query params, don't do anything
			if region := r.URL.Query().Get("region"); region != "" {
				next.ServeHTTP(w, r)
				return
			}

			primary, err := db.GetPrimary()
			if err != nil {
				log.Println("Error getting primary:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If primary is not empty, redirect
			if primary != "" {
				w.Header().Set("fly-replay", "instance="+primary)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
