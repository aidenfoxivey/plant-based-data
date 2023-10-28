package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	server := &http.Server{
		Addr:           ":80",
		Handler:        service(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Serving (💅) on address %s\n", server.Addr)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/javascript"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.Get("/index.mjs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/index.mjs")
	})

	r.Get("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.css")
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	return r
}

func levenshteinDistance(s string, t string) int {
	if s == t {
		return 0
	}

	m := len(s)
	n := len(t)
	d := make([][]int, m)
	distances := make([]int, n*m)

	// Loop over the rows, slicing each row from the front of the remaining pixels slice.
	for i := range d {
		d[i], distances = distances[:n], distances[n:]
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			d[i][j] = 0
		}
	}
	for i := 0; i < n; i++ {
		d[i][0] = i
	}
	for j := 0; j < n; j++ {
		d[0][j] = j
	}

	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			substitutionCost := 1
			if s[i] == t[j] {
				substitutionCost = 0
			}

			d[i][j] = min(min(d[i-1][j]+1,
				d[i][j-1]+1),
				d[i-1][j-1]+substitutionCost)
		}
	}
	return d[m][n]
}
