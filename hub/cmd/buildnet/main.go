package main

import (
	"encoding/json"
	"fmt"
	buildnetv5connect "github.com/buildnethq/buildnet/hub/gen/buildnet/v5/buildnetv5connect"
	"github.com/buildnethq/buildnet/hub/internal/ledger"
	"github.com/buildnethq/buildnet/hub/internal/resources"
	"github.com/buildnethq/buildnet/hub/internal/transport"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// Set by -ldflags in CI during release.
	product = "BuildNET"
	version = "dev"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("%s %s\n", product, version)

	case "hub":
		if len(os.Args) >= 3 && os.Args[2] == "start" {
			hubStart()
			return
		}
		usage()
		os.Exit(2)

	case "workers":
		if len(os.Args) >= 3 && os.Args[2] == "ls" {
			fmt.Println("No workers yet (alpha wedge).")
			return
		}
		usage()
		os.Exit(2)

	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`BuildNET (alpha)

Usage:
  buildnet version
  buildnet hub start
  buildnet workers ls
`)
}

func hubStart() {
	addr := "127.0.0.1:7444"
	mux := http.NewServeMux()

	// --- alpha protocol stubs ---
	transportSvc := transport.New()
	ledgerSvc := ledger.New()
	resourceSvc := resources.New()

	mux.Handle(buildnetv5connect.NewTransportServiceHandler(transportSvc))
	mux.Handle(buildnetv5connect.NewLedgerServiceHandler(ledgerSvc))
	mux.Handle(buildnetv5connect.NewResourceServiceHandler(resourceSvc))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"ok":      true,
			"product": product,
			"version": version,
			"time":    time.Now().UTC().Format(time.RFC3339Nano),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("%s hub starting on http://%s\n", product, addr)
	log.Printf("health: http://%s/healthz\n", addr)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
