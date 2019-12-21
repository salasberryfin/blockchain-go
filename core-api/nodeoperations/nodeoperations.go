package nodeoperations

import (
    "fmt"
    "log"
    "time"
    "errors"
    "net/http"
)

const periodicity time.Duration = 10000

func CheckForInactiveNodes(nodes *[]string) {
    for {
        for _, v := range *nodes {
            _, err := http.Get("http://" + v)
            if err != nil {
                log.Printf("Node %v is not responding -> removing from network!", v)
                newNodes, _ := RemoveNode(*nodes, v)
                (*nodes) = newNodes
                continue
            }
            log.Printf("Node %v is alive.", v)
        }
        time.Sleep(periodicity * time.Millisecond)
    }
}

func ValidateOperation(op string) bool {
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

func AddNode(nodes []string, apiURL string) ([]string, error) {
    for _, b := range nodes {
        if b == apiURL {
            return []string{}, errors.New(fmt.Sprintf("Node %v is already in the network.", apiURL))
        }
    }
    nodes = append(nodes, apiURL)

    return nodes, nil
}

func RemoveNode(nodes []string, apiURL string) ([]string, error) {
    for i, b := range nodes {
        if b == apiURL {
            nodes[i] = nodes[len(nodes) - 1]
            nodes = nodes[:len(nodes) - 1]
            return nodes, nil
        }
    }

    return []string{}, errors.New(fmt.Sprintf("Node %v is not part of the network.", apiURL))
}
