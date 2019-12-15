package main

import (
    "log"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
)

const port string = ":5000"

func convertToJson(parsed []string) ([]byte, error) {
    return json.Marshal(parsed)
}

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    resp := []string{"You're running a valid node in port " + port}
    js, err := convertToJson(resp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Write(js)
}

func NodeApi() {
    r := mux.NewRouter()

    //r.HandleFunc("/new_transaction/", manageNode).Methods(http.MethodPost)
    //r.HandleFunc("/get_nodes", getNodes).Methods(http.MethodGet)
    r.HandleFunc("/", home).Methods(http.MethodGet)

    err := http.ListenAndServe(port, r)
    if err != nil {
        log.Fatal("Server says: ", err)
    }
}
