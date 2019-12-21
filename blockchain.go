package main

import (
    "os"
    "fmt"
    "log"
    "crypto/sha256"
    "io/ioutil"
    "time"
    "strings"
    "strconv"
    "net/http"
    "encoding/hex"
    "encoding/json"
    "errors"
    "bytes"
)

var nodeApiURL string = "localhost"
const mainApi string = "http://localhost:8080"

var Nodes []string
var blockchain []Block;

func getSha256 (unhashed string) string {
    unhashedBytes := []byte(unhashed)
    hashed := sha256.New()
    hashed.Write(unhashedBytes)

    return hex.EncodeToString(hashed.Sum(nil))
}

func validateBlock(block Block) bool {
    if block.PreviousHash != getSha256(fmt.Sprintf("%v", getPreviousBlock())) {
        log.Print("Hashes do not match. Aborting!")
        return false
    }
    hashedProof := getSha256(strconv.Itoa(block.Proof))
    if hashedProof[len(hashedProof)-6:] != "000000" {
        log.Print("Proof of work is wrong. Aborting!")
        return false
    }

    return true
}

func getPreviousBlock() Block {
    if len(blockchain) > 0 {
        return blockchain[len(blockchain)-1]
    }

    return Block{}
}

func generateTransaction (src, rec int, amt string) Transaction {
    var transaction Transaction
    nodesSlice, err := retrieveNodes()
    if err != nil {
        log.Print("[Error]: Error when retrieving the current node network.")
        return Transaction{}
    }
    Nodes = nodesSlice
    nodePort, _ := strconv.Atoi(strings.Split(nodeApiURL, ":")[1])
    if src == nodePort {
        broadcastTransaction(src, rec, amt)
    }
    transaction.Source = src
    transaction.Recipient = rec
    transaction.Amount = amt
    go generateBlock(transaction)

    return transaction
}

func generateBlock (transaction Transaction) (Block, error) {
    var block Block
    nodePort, _ := strconv.Atoi(strings.Split(nodeApiURL, ":")[1])
    timestamp := time.Now().Unix()
    block.Miner = nodePort
    block.Timestamp = timestamp
    block.PreviousHash = getSha256(fmt.Sprintf("%v", getPreviousBlock()))
    block.Transaction = transaction
    proofValue := proofOfWork()
    block.Proof = proofValue
    if validateBlock(block) {
        blockchain = append(blockchain, block)
        updateChainForAllNodes()
        return block, nil
    }

    return Block{}, errors.New("Block was not properly validated.")
}

func broadcastTransaction (src, rec int, amt string) {
    var u string
    for i := 0; i < len(Nodes); i++ {
        u = Nodes[i][1:len(Nodes[i]) - 1]
        transactionUrl := fmt.Sprintf("http://%v/new_transaction/amount/%v/source/%v/recipient/%v", u, amt, src, rec)
        requestBody := []byte{}
        log.Print("Broadcasting transaction: ", transactionUrl)
        _, err := http.Post(transactionUrl, "application/json", bytes.NewBuffer(requestBody))
        if err != nil {
            log.Print("[Error] An error was encountered when broadcasting the new transaction.")
            return
        }
    }
}

func updateChainForAllNodes () {
    var u string
    for i := 0; i < len(Nodes); i++ {
        log.Print("Updating: ", Nodes[i])
        u = Nodes[i][1:len(Nodes[i]) - 1]
        _, err := http.Get("http://" + u + "/update_chain")
        if err != nil {
            log.Print("[Error]: an error was encountered when updating the chain.")
        }
    }
}

func checkForLongerChain() ([]string) {
    nodesSlice, err := retrieveNodes()
    if err != nil {
        log.Print("[Error]: Error when retrieving the current node network.")
        return []string{"An error was encountered when updating the nodes list %v"}
    }
    Nodes = nodesSlice
    for i := 0; i < len(Nodes); i++ {
        var newBlockchain []Block
        u := Nodes[i][1:len(Nodes[i]) - 1]
        log.Print("Sending request to: ", u)
        r, err := http.Get("http://" + u + "/get_chain")
        if err != nil {
            log.Print("[Error]: ", err)
            return []string{fmt.Sprintf("An error was encountered when calling %v", u)}
        }
        defer r.Body.Close()
        responseData, err := ioutil.ReadAll(r.Body)
        json.Unmarshal(responseData, &newBlockchain)
        fmt.Printf("Parsed JSON Response: %v", newBlockchain)
        if len(blockchain) < len(newBlockchain) {
            blockchain = newBlockchain
        }
    }

    return []string{"Chain is now up to date - You can check the current status by calling /get_chain"}
}

func getCurrentChain() (string) {
    var resp string
    for i := 0; i < len(blockchain); i++ {
        b, err := json.Marshal(blockchain[i])
        if err != nil {
            log.Fatal(err)
        }
        resp += string(b)
        if i < len(blockchain) - 1 {
            resp += ","
        }
    }

    return resp
}

func retrieveNodes() ([]string, error) {
    r, err := http.Get(mainApi + "/get_nodes")
    if err != nil {
        return []string{}, err
    }
    defer r.Body.Close()
    responseData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return []string{}, err
    }
    nodesString := strings.Trim(string(responseData), "[]")
    if len(nodesString) == 0 {
        return []string{}, nil
    }
    nodesSlice := strings.Split(nodesString, ",")
    for i := 0; i < len(nodesSlice); i++ {
        u := nodesSlice[i][1:len(nodesSlice[i]) - 1]
        if u == nodeApiURL {
            nodesSlice[i] = nodesSlice[len(nodesSlice) - 1]
            nodesSlice = nodesSlice[:len(nodesSlice) - 1]
        }
    }

    return nodesSlice, nil
}

func main() {
    portArg := os.Args[1:]
    nodeApiPort, err := strconv.Atoi(portArg[0])
    if err != nil {
        log.Fatal("Wrong input arguments.")
    }
    nodeApiURL += ":" + strconv.Itoa(nodeApiPort)
    initNode(nodeApiURL, mainApi)
    initApi(nodeApiPort)
}

