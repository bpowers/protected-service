// Copyright 2013 Bobby Powers. All rights reserved.

package main

import (
	"net/http"
	"time"
)

type decacheHandler struct {
	http.Handler
}

type cacheHandler struct {
	http.Handler
}

var (
	past   = today().Add(-24 * time.Hour)     // yesterday
	future = today().Add(60 * 24 * time.Hour) // 6 months from now
)

// midnight at the start of today
func today() time.Time {
	now := time.Now().Unix()
	return time.Unix(now-now%int64(24*60*60), 0)
}

func addExpiresHeaderFor(path string) bool {
	return true
}

func forceRevalidate(w http.ResponseWriter) {
	w.Header().Set("Expires", past.Format(http.TimeFormat))
}

func forceCache(w http.ResponseWriter) {
	w.Header().Set("Expires", future.Format(http.TimeFormat))
}

// by setting the expires tag to yesterday, we ensure that browsers
// _always_ do a GET for resources with an If-Modified-Since header.
// If they have the latest copy they get a 304 response and use their
// cache, but the point is they always check.
func (h *decacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if addExpiresHeaderFor(r.URL.Path) {
		forceRevalidate(w)
	}
	h.Handler.ServeHTTP(w, r)
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forceCache(w)
	h.Handler.ServeHTTP(w, r)
}
