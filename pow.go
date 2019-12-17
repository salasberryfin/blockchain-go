package main

import (
    "log"
    "math/rand"
    "strconv"
)


func proofOfWork () bool {
    var x, y int
    var hashed string
    for {
        x = rand.Intn(10000)
        y = rand.Intn(10000)
        hashed = getSha256(strconv.Itoa(x * y))
        log.Print(hashed)
        if hashed[len(hashed)-4:] == "0000" {
            log.Printf("Got it for x=%d, y=%d.", x, y)
            break
        }
    }

    return true
}
