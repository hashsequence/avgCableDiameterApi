package main

import (
	"fmt"
	"time"
)


func main() {
	t1 := time.Now()
	time.Sleep(time.Duration(2) * time.Second)
	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(diff)
}