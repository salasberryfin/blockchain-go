package main

import (
    "fmt"
    "log"
    "strconv"
    "net/http"
    "io/ioutil"
    "strings"
    "bytes"
)

type Block struct {
    Timestamp       int64           `json:"Timestamp"`
    Miner           int             `json:"Miner"`
    Proof           int             `json:"Proof"`
    PreviousHash    string          `json:"PreviousHash"`
    Transaction     Transaction     `json:"Transaction"`
}

type Transaction struct {
    Source      int
    Recipient   int
    Amount      string
}

func initNode (nodeUrl, mainUrl string) {
    registerUrl := fmt.Sprintf("%v/manage_node/operation/add/%v", mainUrl, nodeUrl)
    requestBody := []byte{}
    r, err := http.Post(registerUrl, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        log.Print("[Error] An error was encountered when registering the node.")
        return
    }
    defer r.Body.Close()
    responseData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Print("[Error] An error was encountered when parsing the server response.")
        return
    }
    responseString := strings.Trim(string(responseData), "[]")
    log.Print(responseString)
}

func initApi(port int) {
    checkForLongerChain()
    log.Printf("Node API is running - port: %v", port)
    NodeApi(":" + strconv.Itoa(port))
}

