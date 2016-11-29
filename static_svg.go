// Copyright 2013 Bobby Powers. All rights reserved.

package main

import (
	"net/http"
	"strings"
)

/// svgHandler fixes the Content-Type response header for SVG files.
/// Go's built in content-sniffer incorrectly sets the content type as
/// 'text/xml', which results in browsers not displaying SVGs as
/// images.
type svgHandler struct {
	http.Handler
}

func (h *svgHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, ".svg") {
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	h.Handler.ServeHTTP(w, req)
}
