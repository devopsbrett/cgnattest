package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type portFrequency struct {
	mu    sync.Mutex
	ports map[string]int
}

func (p *portFrequency) portUsed(port string) {
	p.mu.Lock()
	p.ports[port]++
	p.mu.Unlock()
}

func (p *portFrequency) results() {
	reused := 0
	for _, v := range p.ports {
		if v > 1 {
			reused += 1
		}
	}
	fmt.Printf("Reused ports: %d/%d (%0.4f%%)\n", reused, len(p.ports), (float64(reused)/float64(len(p.ports)))*100.0)
}

// var workerChan chan bool
var wg sync.WaitGroup
var pf portFrequency

func dialServer(wc <-chan bool, serverURL string) {
	defer wg.Done()

	for _ = range wc {
		c, err := net.Dial("tcp4", serverURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg, _ := bufio.NewReader(c).ReadString('\n')
		addrParts := strings.SplitN(msg, ":", 2)
		pf.portUsed(addrParts[1])
		c.Close()
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("syntax: %s server_uri [connection_count]\n", os.Args[0])
		fmt.Println("server_uri is mandatory and should be given in format hostip:hostport")
		fmt.Println("connection_count is optional and will default to 2048")
		return
	}

	pf = portFrequency{ports: make(map[string]int)}

	workerChan := make(chan bool)

	workerCount := runtime.NumCPU()

	connectionCount := 2048
	// portFrequency := make(map[int]int)
	connect := os.Args[1]

	if len(os.Args) > 2 {
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid connection_count given")
			return
		}
		connectionCount = num
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go dialServer(workerChan, connect)
	}

	for i := 0; i < connectionCount; i++ {
		workerChan <- true
	}

	close(workerChan)

	wg.Wait()
	pf.results()
}
