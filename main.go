// Copyright 2016 Bobby Powers. All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/bpowers/seshcookie"
)

var (
	devMode bool
	devUrl  string
)

func logGzipAndCORS(h http.Handler) http.Handler {
	return &loggedHandler{&corsHandler{h}}
}

var serveCurrentUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	details := make(map[string]string)

	session := seshcookie.Session.Get(r)
	details["username"], _ = session["user"].(string)

	detailsBytes, err := json.Marshal(details)
	if err != nil {
		log.Printf("json.Marshal(%#v): %s", details, err)
		return
	}

	w.Write(detailsBytes)
})

func main() {
	addr := flag.String("addr", "127.0.0.1:8009", "address to listen on")
	flag.BoolVar(&devMode, "dev", false, "development")
	flag.StringVar(&devUrl, "dev-url", "http://127.0.0.1:8080", "URL to proxy for development")
	flag.Parse()

	config, err := ReadConfig("config")
	if err != nil {
		log.Printf("ReadConfig(): %s", err)
		return
	}

	var staticHandler http.Handler
	staticHandler = http.FileServer(http.Dir("./static"))

	// XXX: for development - we proxy html/css/js requrests to
	// grunt's server task.  Note that devMode is invasive and has
	// its hooks in auth and manager too
	if devMode {
		devServer, err := url.Parse(devUrl)
		if err != nil {
			log.Fatalf("url.Parse: %s", err)
		}
		p := httputil.NewSingleHostReverseProxy(devServer)
		staticHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//r.URL, _ = url.Parse("/")
			p.ServeHTTP(w, r)
		})
	}

	rootHandler := seshcookie.NewSessionHandler(
		&AuthHandler{
			&decacheHandler{&svgHandler{staticHandler}},
			&authorizer{config.AuthDir},
			&decider{},
			true,
		},
		config.SessionKey,
		nil)
	rootHandler.CookieName = config.CookieName

	currentUserHandler := seshcookie.NewSessionHandler(
		&AuthHandler{
			&decacheHandler{serveCurrentUser},
			&authorizer{config.AuthDir},
			&decider{},
			true,
		},
		config.SessionKey,
		nil)
	currentUserHandler.CookieName = config.CookieName

	http.Handle("/", logGzipAndCORS(rootHandler))
	http.Handle("/api/user/current", logGzipAndCORS(currentUserHandler))
	http.Handle("/err/",
		&decacheHandler{
			http.StripPrefix("/err/",
				http.FileServer(http.Dir("./err")))})

	log.Printf("listening on %s", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Printf("ListenAndServe: %s", err)
	}
}
