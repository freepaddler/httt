package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ident       = "  "
	defaultPort = "8080"
)

func main() {
	withBody := os.Getenv("WITH_BODY") != ""
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	if host == "" {
		host, _ = os.Hostname()
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("server stated at port %s", port)
	if err := http.ListenAndServe(":"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RemoteAddr:%s, Host:%s, Method:%s, URI:%s", r.RemoteAddr, r.Host, r.Method, r.URL)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Httt-Host", host)

		_, _ = fmt.Fprintf(w, "Headers:\n")
		for name, values := range r.Header {
			for _, value := range values {
				_, _ = fmt.Fprintf(w, "%s%s: %s\n", ident, name, value)
			}
		}

		_, _ = fmt.Fprintln(w)

		_, _ = fmt.Fprintf(w, "Proto: %s\n", r.Proto)
		_, _ = fmt.Fprintf(w, "Content-Length: %d\n", r.ContentLength)
		_, _ = fmt.Fprintf(w, "Transfer-Encoding: %s\n", r.TransferEncoding)

		_, _ = fmt.Fprintln(w)

		_, _ = fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
		if authHeader := r.Header.Get("Authorization"); authHeader != "" && strings.HasPrefix(authHeader, "Basic ") {
			authHeader = strings.TrimPrefix(authHeader, "Basic ")
			if creds, err := base64.StdEncoding.DecodeString(authHeader); err == nil {
				_, _ = fmt.Fprintf(w, "Credentials: %s\n", creds)
			}
		}
		_, _ = fmt.Fprintln(w)

		_, _ = fmt.Fprintf(w, "Host: %s\n", r.Host)
		_, _ = fmt.Fprintf(w, "Method: %s\n", r.Method)
		_, _ = fmt.Fprintf(w, "URL: %s\n", r.RequestURI)

		if withBody {
			if body, err := io.ReadAll(r.Body); err == nil {
				_, _ = fmt.Fprintf(w, "\nBody:\n")
				_, _ = fmt.Fprintf(w, "%s", body)
			}
		}

	})); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
