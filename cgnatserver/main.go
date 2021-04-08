package main

import (
	"fmt"
	"net"
	"os"
)

func handleConnection(c net.Conn) {
	fmt.Print(".")
	remoteAddr := c.RemoteAddr().String()
	c.Write([]byte(remoteAddr))
	c.Close()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter a port number")
		return
	}
	port := ":" + os.Args[1]
	l, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
