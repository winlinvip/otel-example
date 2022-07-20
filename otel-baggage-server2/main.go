package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// See https://opentelemetry.io/docs/instrumentation/go/libraries/
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Client header", r.Header)
		io.WriteString(w, "Hello")
	}))

	fmt.Println("Please test by http://localhost:8096")
	http.ListenAndServe(":8096", nil)
}
