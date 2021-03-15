package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartApi() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeApi)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func HomeApi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("wizechain api v0.0.1"))
}

func JSONResponse(w http.ResponseWriter, code int, output interface{}) {
	// Convert our interface to JSON
	response, _ := json.Marshal(output)
	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")
	// Our response code
	w.WriteHeader(code)

	w.Write(response)
}
