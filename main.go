package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	fmt.Printf("Server initiated")
	http.ListenAndServe(":3000", r)
}
