package main

import "fmt"
import "time"

func main() {
    start := time.Now()
    fetch("http://www.google.com")
    fmt.Println(time.Since(start))
}
