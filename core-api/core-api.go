package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/salasberryfin/blockchain-go/core-api/nodeoperations"
    "github.com/gorilla/mux"
)

var Nodes = []string{}

func manageNode(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    w.Header().Set("Content-Type", "application/json")
    oper, api_url := "", ""
    if val_api, valid := params["api_url"]; valid {
        api_url = val_api
    }
    if val_oper, valid := params["oper"]; valid && nodeoperations.ValidateOperation(val_oper) {
        oper = val_oper
    } else {
        resp := []string{"Invalid operation '" + val_oper + "' was requested."}
        js, err := json.Marshal(resp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(js)

        return
    }

    switch oper {
    case "add":
        resp, err := nodeoperations.AddNode(Nodes, api_url)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        Nodes = resp
        js, err := json.Marshal([]string{fmt.Sprintf("%v is now part of the blockchain network.", api_url)})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write(js)

        return
    case "remove":
        resp, err := nodeoperations.RemoveNode(Nodes, api_url)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        Nodes = resp
        js, err := json.Marshal([]string{fmt.Sprintf("%v has been removed from the blockchain network.", api_url)})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write(js)

        return
    }
}

func getNodes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    js, err := json.Marshal(Nodes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Write(js)
}

func main() {
    log.Print("Server running on port 8080...")
    r := mux.NewRouter()
    r.HandleFunc("/manage_node/operation/{oper}/{api_url}", manageNode).Methods(http.MethodPost)
    r.HandleFunc("/get_nodes", getNodes).Methods(http.MethodGet)
    go nodeoperations.CheckForInactiveNodes(&Nodes)
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Server says: ", err)
    }
}
