package main

import (
    "log"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
)

var Nodes = []string{}

func convertToJson(parsed []string) ([]byte, error) {
    return json.Marshal(parsed)
}

func validateOperation(op string) bool {
    operations := []string{"add", "remove"}
    for _, b := range operations {
        if b == op {
            log.Print("Detected valid operation: ", b)
            return true
        }
    }
    log.Print("Detected invalid operation: ", op)

    return false
}

func manageNode(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    w.Header().Set("Content-Type", "application/json")
    oper, api_url := "", ""
    if val_api, valid := params["api_url"]; valid {
        api_url = val_api
    }
    if val_oper, valid := params["oper"]; valid && validateOperation(val_oper) {
        oper = val_oper
    } else {
        resp := []string{"Invalid operation '" + val_oper + "' was requested."}
        js, err := convertToJson(resp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(js)
        return
    }

    if oper == "add" {
        for i := 0; i < len(Nodes); i++ {
            if Nodes[i] == api_url {
                resp := []string{"Node is already in the list."}
                js, err := convertToJson(resp)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                w.Write(js)
                return
            }
        }
        Nodes = append(Nodes, api_url)
        resp := []string{"Adding node: " + api_url}
        js, err := convertToJson(resp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(js)
    } else if oper == "remove" {
        for i, b := range Nodes {
            if b == api_url {
                Nodes[len(Nodes)-1], Nodes[i] = Nodes[i], Nodes[len(Nodes)-1]
                Nodes = Nodes[:len(Nodes)-1]
                resp := []string{"Node with API '" + api_url + "' was removed from list."}
                js, err := convertToJson(resp)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                w.Write(js)
                return
            }
        }
        resp := []string{"Node with API '" + api_url + "' is not part of the network."}
        js, err := convertToJson(resp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(js)
    }
}

func getNodes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    js, err := convertToJson(Nodes)
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
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Server says: ", err)
    }
}
