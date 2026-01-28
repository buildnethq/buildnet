package main

import (
	"encoding/json"
	"fmt"
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
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("%s hub starting on http://%s\n", product, addr)
	log.Printf("health: http://%s/healthz\n", addr)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
