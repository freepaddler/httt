package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ident       = "  "
	defaultPort = "8080"
)

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total_labeled",
		Help: "Total number of HTTP requests, labeled by label",
	},
	[]string{"label"},
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

	prometheus.MustRegister(requestCounter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("server stated at port %s", port)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get("X-Request-Id")
		if reqId == "" {
			reqId = uuid.NewString()
		}

		label := r.Header.Get("label")
		fmt.Println("label:", label)
		requestCounter.WithLabelValues(label).Inc()

		var buf bytes.Buffer
		lw := io.Writer(&buf)
		_, _ = fmt.Fprintf(lw, "Request ID: %s\n", reqId)
		m := io.MultiWriter(w, lw)

		defer r.Body.Close()

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Httt-Host", host)
		w.Header().Set("X-Request-Id", reqId)

		_, _ = fmt.Fprintf(m, "Headers:\n")
		for name, values := range r.Header {
			for _, value := range values {
				_, _ = fmt.Fprintf(m, "%s%s: %s\n", ident, name, value)
			}
		}

		_, _ = fmt.Fprintf(m, "Proto: %s\n", r.Proto)
		_, _ = fmt.Fprintf(m, "Content-Length: %d\n", r.ContentLength)
		_, _ = fmt.Fprintf(m, "Transfer-Encoding: %s\n", r.TransferEncoding)

		_, _ = fmt.Fprintf(m, "RemoteAddr: %s\n", r.RemoteAddr)
		if authHeader := r.Header.Get("Authorization"); authHeader != "" && strings.HasPrefix(authHeader, "Basic ") {
			authHeader = strings.TrimPrefix(authHeader, "Basic ")
			if creds, err := base64.StdEncoding.DecodeString(authHeader); err == nil {
				_, _ = fmt.Fprintf(m, "Credentials: %s\n", creds)
			}
		}

		_, _ = fmt.Fprintf(m, "Host: %s\n", r.Host)
		_, _ = fmt.Fprintf(m, "Method: %s\n", r.Method)
		_, _ = fmt.Fprintf(m, "URL: %s\n", r.RequestURI)

		if withBody {
			if body, err := io.ReadAll(r.Body); err == nil && len(body) > 0 {
				_, _ = fmt.Fprintf(w, "Body:\n")
				_, _ = fmt.Fprintf(w, "%s", body)
			}
		}
		log.Println(buf.String())
	})

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":"+port, nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
