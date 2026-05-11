package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Logger adalah middleware untuk mencatat log request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Lanjutkan ke handler berikutnya (seperti next() di Express)
		next.ServeHTTP(w, r)

		// Setelah handler selesai, catat durasinya
		fmt.Printf(
			"[%s] %s %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
