package main

import (
    "fmt"
    "log"
    "crypto/sha256"
    "hash"
    "io/ioutil"
    //"errors"
    "net/http"
    //"encoding/json"
)

const main_api string = "http://localhost:8080/get_nodes"
var Nodes = []string{}

type Block struct {
    Hash hash.Hash
    PreviousHash hash.Hash
    Transaction string
}

var blockchain []Block;

func getSha256 (unhashed string) hash.Hash {
    unhashedBytes := []byte(unhashed)
    hashed := sha256.New()
    hashed.Write(unhashedBytes)

    return hashed
}

func getPreviousBlock() Block {
    block := Block{}

    return block
}

func validateBlock() bool {

    return false
}

func addToChain(block Block) {
}

func generateBlock (hash hash.Hash, transaction string) Block {
    block := Block{}
    block.Hash = hash
    block.PreviousHash = getSha256(fmt.Sprintf("%v", getPreviousBlock()))
    block.Transaction = transaction

    return block
}

func retrieveNodes() (string, error) {
    r, err := http.Get(main_api)
    if err != nil {
        return "", err
    }
    defer r.Body.Close()
    responseData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    responseString := string(responseData)
    log.Print(responseString)

    return responseString, nil
}

func main() {
    //hash := getSha256("new block")
    //log.Printf("%x", hash.Sum(nil))
    //transaction := "transaction data"
    //new_block := generateBlock(hash, transaction)
    //log.Print("new block: ", new_block.Transaction)
    log.Print("Starting API...")
    Nodes, err := retrieveNodes()
    if err != nil {
        log.Fatal("Could not retrieve the current cluster.")
    }
    log.Print("Updated list of nodes: ", Nodes)
    NodeApi()
}

