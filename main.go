package main

import (
    "fmt"
    "log"
    "crypto/sha256"
    "io/ioutil"
    "time"
    "strings"
    "net/http"
    "encoding/hex"
    "encoding/json"
)

const node_api_port string = ":5000"
const node_api_url string = "localhost" + node_api_port
const main_api string = "http://localhost:8080/get_nodes"
var Nodes = []string{}

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

func validateBlock() bool {

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

func generateBlock (transaction Transaction) Block {
    block := Block{}
    timestamp := time.Now().Unix()
    block.Proof = "something"
    block.Timestamp = timestamp
    block.PreviousHash = getSha256(fmt.Sprintf("%v", getPreviousBlock()))
    block.Transaction = transaction

    return block
}


func checkForLongerChain() ([]string) {
    nodesString, err := retrieveNodes()
    Nodes = strings.Split(nodesString, ",")
    if err != nil {
        log.Fatal(err)
    }
    for i := 0; i < len(Nodes); i++ {
        u := Nodes[i]
        u = u[1:len(u)-1]
        if u == node_api_url {
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
        var newBlockchain []Block
        json.Unmarshal(responseData, &newBlockchain)
        fmt.Printf("Parsed JSON Response: %v", newBlockchain)
        if len(blockchain) < len(newBlockchain) {
            log.Print("Retrieved blockchain is longer.")
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
    r, err := http.Get(main_api)
    if err != nil {
        return "", err
    }
    defer r.Body.Close()
    responseData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    responseString := strings.Trim(string(responseData), "[]")

    return responseString, nil
}

func main() {
    transaction := Transaction{0, 0, "Genesis"}
    new_block := generateBlock(transaction)
    addToChain(new_block)
    transaction2 := Transaction{1, 1, "Second"}
    new_block2 := generateBlock(transaction2)
    addToChain(new_block2)
    log.Printf("Node API is running - port: %v", strings.Trim(node_api_port, ":"))
    NodeApi(node_api_port)
}

