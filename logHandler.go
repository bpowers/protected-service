// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogMsg struct {
	Path         string  `json:"path"`
	Referer      string  `json:"referer"`
	Method       string  `json:"method"`
	UserAgent    string  `json:"user_agent"`
	IpAddress    string  `json:"ip_address"`
	TimeUnix     int64   `json:"time"`
	Duration     float64 `json:"duration"`
	ResponseCode int32   `json:"response_code,omitempty"`
}

// LoggedHandler wraps a http.Handler, writing log information in JSON
// form to stdout.
type loggedHandler struct {
	http.Handler
}

func (h *loggedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ipAddr := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ipAddr == "" && r.RemoteAddr != "" {
		ipAddr = r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	}

	start := time.Now()
	h.Handler.ServeHTTP(w, r)
	duration := time.Now().Sub(start)

	msg := &LogMsg{
		Path:      r.URL.Path,
		Referer:   r.Header.Get("Referer"),
		UserAgent: r.Header.Get("User-Agent"),
		Method:    r.Method,
		IpAddress: ipAddr,
		TimeUnix:  start.Unix(),
		Duration:  duration.Seconds(),
		//ResponseCode: w.StatusCode,
	}
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("json.Marshal(%#v): %s", msg, err)
	} else {
		msgBuf = append(msgBuf, '\n')
		os.Stdout.Write(msgBuf)
	}
}
