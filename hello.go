package main

import (
	"fmt"
	"time"
)

func main() {
	mymap := make(map[int]int)

	start := time.Now()
	for i := 0; i < 100666999; i++ {
		mymap[i] = i
	}
	lenghth := len(mymap)
	elapsed := time.Since(start)
	fmt.Println(elapsed, lenghth)
}
