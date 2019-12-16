package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
)

func convertToJson(parsed []string) ([]byte, error) {
    return json.Marshal(parsed)
}

func updateChain(w http.ResponseWriter, r *http.Request) {
    resp := checkForLongerChain()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    js, err := convertToJson(resp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Write(js)
}

func getChain(w http.ResponseWriter, r *http.Request) {
    resp := getCurrentChain()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    for i := 0; i < len(resp); i++ {
    }
    js := []byte(fmt.Sprintf("[%v]", resp))
    w.Write(js)
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
    /**
    Generate new transaction:
    - source
    - destination
    - amount
    Identify nodes by port in which the API is running
    **/
    //source := strings.Split(port, ":")[1]
    //params := mux.Vars(r)
    //w.Header().Set("Content-Type", "application/json")
}

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    resp := []string{"You're running a valid node"}
    js, err := convertToJson(resp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Write(js)
}

func NodeApi(port string) {
    r := mux.NewRouter()
    r.HandleFunc("/", home).Methods(http.MethodGet)
    r.HandleFunc("/new_transaction/amount/{amount}/recipient/{target}", newTransaction).Methods(http.MethodPost)
    r.HandleFunc("/get_chain", getChain).Methods(http.MethodGet)
    r.HandleFunc("/update_chain", updateChain).Methods(http.MethodGet)

    err := http.ListenAndServe(port, r)
    if err != nil {
        log.Fatal("Server says: ", err)
    }
}
