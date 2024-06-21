package main

// import (
// 	"log"
// 	"net/http"
// )

// func middlewareLogger(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.path)
// 		next.ServeHTTP(w, r)
// 	})
// }
