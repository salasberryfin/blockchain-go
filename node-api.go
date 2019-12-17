package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"
    "strings"
    "strconv"

    "github.com/gorilla/mux"
)

var apiPort string

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
    // Generate new Transaction{Source, Destination, Amount}
    // Identify nodes (Source, Destination) by API port
    w.Header().Set("Content-Type", "application/json")
    var amount string
    var source, recipient int
    params := mux.Vars(r)
    source, err := strconv.Atoi(strings.Split(apiPort, ":")[1])
    if err != nil {
        js, _ := convertToJson([]string{"Failed to parse request parameters."})
        w.Write(js)
        return
    }
    if val_amount, valid := params["amount"]; valid {
        amount = val_amount
    }
    if val_recipient, valid := params["target"]; valid {
        recipient, err = strconv.Atoi(val_recipient)
        if err != nil {
            js, _ := convertToJson([]string{"Failed to parse request parameters."})
            w.Write(js)
            return
        }
    }
    var block Block
    block, err = generateBlock(source, recipient, amount)
    if err != nil {
        log.Print("Error found")
    }
    log.Print(block)
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
    apiPort = port
    r := mux.NewRouter()
    r.HandleFunc("/", home).Methods(http.MethodGet)
    r.HandleFunc("/new_transaction/amount/{amount}/recipient/{target}", newTransaction).Methods(http.MethodPost)
    r.HandleFunc("/get_chain", getChain).Methods(http.MethodGet)
    r.HandleFunc("/update_chain", updateChain).Methods(http.MethodGet)

    err := http.ListenAndServe(apiPort, r)
    if err != nil {
        log.Fatal("Server says: ", err)
    }
}
