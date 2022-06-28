package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

/**
Please note Start functions is a placeholder for you to start your own solution.
Feel free to drop gorilla.mux if you want and use any other solution available.

main function reads host/port from env just for an example, flavor it following your taste
*/

// Start /** Starts the web server listener on given host and port.
func Start(host string, port int) {
	router := mux.NewRouter()

	// Register routes (endpoints)
	router.HandleFunc("/name/{PARAM}", HomeHandler).Methods(http.MethodGet)
	router.HandleFunc("/bad", BadHandler).Methods(http.MethodGet)
	router.HandleFunc("/data", DataHandler).Methods(http.MethodPost)
	router.HandleFunc("/headers", HeaderHandler).Methods(http.MethodPost)

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), loggingMiddleware(router)); err != nil {
		log.Fatal(err)
	}
}

// Logging incomming requests, male lower case and remove trailing slashes
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		r.URL.Path = strings.ToLower(r.URL.Path)
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Return home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello, %v!", params["PARAM"])))
}

// Return bad request
func BadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// Return in body posted content from request body
func DataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte("I got message:\n"))
	w.Write(body)
}

// Return in headers posted content from request headers
func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.Atoi(r.Header.Get("a"))
	if err != nil {
		log.Println(err)
	}
	b, err := strconv.Atoi(r.Header.Get("b"))
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("a+b", fmt.Sprintf("\"%v\"", a+b))
	w.WriteHeader(http.StatusOK)
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
