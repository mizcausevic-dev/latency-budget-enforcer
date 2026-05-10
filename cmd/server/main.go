package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/cli"
	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/httpapi"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "-mode" && len(os.Args) > 2 && os.Args[2] == "cli" {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--mode=cli" {
		return
	}
	if len(os.Args) > 2 && os.Args[1] == "-mode" && os.Args[2] == "cli" {
		return
	}

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: httpapi.NewServer(),
	}

	log.Printf("latency-budget-enforcer listening on http://%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
