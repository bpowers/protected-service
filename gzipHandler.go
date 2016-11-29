// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipHandler struct {
	http.Handler
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (f *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		w = gzipResponseWriter{Writer: gz, ResponseWriter: w}
	}

	f.Handler.ServeHTTP(w, r)
}
