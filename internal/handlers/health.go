package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler diawali huruf Kapital agar bisa di-import oleh package main (Public)
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
