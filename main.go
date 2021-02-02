package main

import (
	"fmt"
	"net"
	"sort"
)

// 1 - The worker(ports, results chan int) function has been modified to accept two channels
func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("192.168.1.219:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// 2 - If the port is closed, channel a 0
			results <- 0
			continue
		}
		conn.Close()
		// 3 - If the port is open, channel to results
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	// 4 - Create a separate channel to communicate the results from the worker to the main thread
	results := make(chan int)
	// 5 - Create Slice to store results for later
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// 6 - Send to the workers in a separate goroutine because the result-gathering loop needs to start before more than 100 items of work can continue
	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
	}()

	// 7 - The result-gathering loop receives on the "results" channel 65535 times
	for i := 0; i < 65535; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	// 8 - After closing the channels, you'll use sort to sort the slice of open ports
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
