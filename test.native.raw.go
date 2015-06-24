package main

import (
	"fmt"
	"os"
	"runtime"
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

	// Building holder
	box := noncense.NewNoncesAdderNative(count);

	fmt.Printf("Built map with %v items\n\n", count)
	for i := 0; i < count; i++ {
		box.AddSync(fmt.Sprintf("%v", i));
		if i % 50000 == 0 {
			fmt.Printf("Done %v entries\n", i);
		}
	}

	fmt.Println("Finished")
}