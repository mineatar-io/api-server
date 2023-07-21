package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
)

// Profile is a Goroutine for automatically profiling the program for optimization and debug reasons.
func Profile(port uint16) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pprof.Profile)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
