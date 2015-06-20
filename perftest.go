package main

import (
	"fmt"
	"runtime"
	"./lib"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Building listener
	count := 1000000
	box := noncense.NewNoncesAdder(uint32(count))

	stop := make(chan bool)

	fmt.Printf("Built map with %v items\n\n", count)
	for i := 0; i < count; i++ {
		go func(x int) {
			res := <- box.Add(fmt.Sprintf("%v", x))
			if x % 10000 == 0 {
				fmt.Printf("Result for %v is %v\n", x, res);
			}
			if x == count - 1 {
				stop <- true
				fmt.Println("DONE");
			}
		}(i)
	}

	_ = <-stop
	fmt.Println("finished")
}