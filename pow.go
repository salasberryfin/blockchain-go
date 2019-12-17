package main

import (
    "log"
    "math/rand"
    "strconv"
)


func proofOfWork () (int) {
    var x, y int
    var hashed string
    log.Print("Executing POW...")
    for {
        x = rand.Intn(10000)
        y = rand.Intn(10000)
        hashed = getSha256(strconv.Itoa(x * y))
        if hashed[len(hashed)-6:] == "000000" {
            log.Printf("Got it for x=%d, y=%d.", x, y)
            break
        }
    }

    return x * y
}
