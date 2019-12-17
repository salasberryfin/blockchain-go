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
)

var nodeApiUrl string = "localhost"
const mainApi string = "http://localhost:8080/get_nodes"
var Nodes []string

type Block struct {
    Timestamp       int64           `json:"Timestamp"`
    Proof           string          `json:"Proof"`
    PreviousHash    string          `json:"PreviousHash"`
    Transaction     Transaction     `json:"Transaction"`
}

type Transaction struct {
    Source      int
    Recipient   int
    Amount      string
}

var blockchain []Block;

func getSha256 (unhashed string) string {
    unhashedBytes := []byte(unhashed)
    hashed := sha256.New()
    hashed.Write(unhashedBytes)

    return hex.EncodeToString(hashed.Sum(nil))
}

func validateBlock(block Block) bool {

    return true
}

func addToChain(block Block) {
    blockchain = append(blockchain, block)
}

func getPreviousBlock() Block {
    if len(blockchain) > 0 {
        return blockchain[len(blockchain)-1]
    }

    return Block{}
}

func generateTransaction (src, rec int, amt string) Transaction {
    var transaction Transaction
    transaction.Source = src
    transaction.Recipient = rec
    transaction.Amount = amt

    return transaction
}

func generateBlock (src, rec int, amt string) (Block, error) {
    var block Block
    var transaction Transaction
    transaction = generateTransaction(src, rec, amt)
    timestamp := time.Now().Unix()
    block.Proof = "proof-tbd"
    block.Timestamp = timestamp
    block.PreviousHash = getSha256(fmt.Sprintf("%v", getPreviousBlock()))
    block.Transaction = transaction
    if validateBlock(block) {
        addToChain(block)
        return block, nil
    }

    return Block{}, errors.New("Block was not properly validated.")
}

func checkForLongerChain() ([]string) {
    nodesString, err := retrieveNodes()
    if err != nil {
        log.Fatal(err)
    }
    if len(nodesString) == 0 {
        return []string{"[Error] No nodes are part of the network."}
    }
    Nodes = strings.Split(nodesString, ",")
    for i := 0; i < len(Nodes); i++ {
        var newBlockchain []Block
        u := Nodes[i][1:len(Nodes[i]) - 1]
        if u == nodeApiUrl {
            log.Print("Avoiding sending a request to ourselves.")
            continue
        }
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

func retrieveNodes() (string, error) {
    r, err := http.Get(mainApi)
    if err != nil {
        return "", err
    }
    defer r.Body.Close()
    responseData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return "", err
    }
    responseString := strings.Trim(string(responseData), "[]")

    return responseString, nil
}

func main() {
    portArg := os.Args[1:]
    nodeApiPort, err := strconv.Atoi(portArg[0])
    if err != nil {
        log.Fatal("Wrong input arguments.")
    }
    nodeApiUrl += ":" + strconv.Itoa(nodeApiPort)
    _, err = generateBlock(0, 0, "Genesis")
    if err != nil {
        log.Print("An error occured when generating a random block.")
    }
    _, err = generateBlock(1, 1, "Second")
    if err != nil {
        log.Print("An error occured when generating a random2 block.")
    }
    log.Printf("Node API is running - port: %v", nodeApiPort)
    proofOfWork()
    NodeApi(":" + strconv.Itoa(nodeApiPort))
}

