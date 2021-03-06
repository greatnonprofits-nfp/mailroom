package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(ww, r)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		elapsed := time.Since(start)
		uri := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
		ww.Header().Set("X-Elapsed-NS", strconv.FormatInt(int64(elapsed), 10))

		if r.RequestURI != "/" {
			log.WithFields(log.Fields{
				"method":     r.Method,
				"status":     ww.Status(),
				"elapsed":    elapsed,
				"length":     ww.BytesWritten(),
				"url":        uri,
				"user_agent": r.UserAgent(),
			}).Info("request completed")
		}
	})
}

// recovers from panics, logs them to sentry and returns an HTTP 500 response
func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				debug.PrintStack()
				log.WithError(errors.New(fmt.Sprint(rvr))).Error("recovered from panic in web handling")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
