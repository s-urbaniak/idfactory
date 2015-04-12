// idfactory is a http service for generating signed UUIDs
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/s-urbaniak/idfactory/signed"
)

var secret []byte

func validate(w http.ResponseWriter, r *http.Request) {
	s, ok := mux.Vars(r)["signed"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	actual, err := signed.Parse(s)
	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	if !actual.Validate(secret) {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func sign(w http.ResponseWriter, r *http.Request) {
	s := signed.New(secret)
	var response struct {
		ID     string `json:"id"`
		Signed string `json:"signed"`
	}

	response.ID = s.ID()
	response.Signed = s.String()

	j, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/%s", s))
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func main() {
	var secretArg string
	var addr string

	flag.StringVar(
		&addr,
		"addr",
		":8080",
		`HTTP service address (e.g., ':6060')`,
	)

	flag.StringVar(
		&secretArg,
		"secret",
		"",
		"The secret to be used when signing generated UUIDs.",
	)

	flag.Parse()

	if secretArg == "" {
		flag.Usage()
		os.Exit(1)
	}

	secret = []byte(secretArg)

	r := mux.NewRouter()

	r.HandleFunc("/{signed:[a-zA-Z0-9-=:/+]+}", validate).
		Methods("GET")

	r.HandleFunc("/", sign).
		Methods("POST")

	log.Fatal(http.ListenAndServe(addr, r))
}
