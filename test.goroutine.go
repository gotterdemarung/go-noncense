package main

import (
	"fmt"
	"runtime"
	"os"
	"strconv"
	"./go-noncense"
)

func main() {
	args  := os.Args[1:]
	cpu   := runtime.NumCPU()
	procs := cpu
	count := 1000000

	if len(args) > 0 {
		i, err := strconv.Atoi(args[0])
		if err != nil || i < 1 {
			fmt.Printf("Not valid process count\n\n")
			os.Exit(2)
		}
		procs = i
	}

	fmt.Println("No concurrency test")
	fmt.Printf(" Processors:   %v\n", cpu)
	fmt.Printf(" Processes:    %v\n", procs)
	fmt.Printf(" Map size:     %v\n", count)

	runtime.GOMAXPROCS(procs)

	// Building listener
	box := noncense.NewNoncesAdder(uint32(count))

	stop := make(chan bool)

	fmt.Printf("Built map with %v items\n\n", count)
	for i := 0; i < count; i++ {
		go func(x int) {
			_ = <- box.Add(fmt.Sprintf("%v", x))
			if x % 50000 == 0 {
				fmt.Printf("Done %v entries\n", x);
			}
			if x == count - 1 {
				stop <- true
			}
		}(i)
	}

	_ = <-stop
	fmt.Println("Finished")
}