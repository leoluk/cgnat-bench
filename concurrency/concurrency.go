package main

import (
	"flag"
	"log"
	"net"
	"time"
)

var (
	addr = flag.String("addr", "[::1]:8080", "Server address")
	n    = flag.Uint("n", 10, "Number of concurrent connections to attempt")
)

func init() {
	flag.Parse()
}

func connect(n uint, done chan struct{}) {
	c, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Printf("connection %d dial failed: %v", n, err)
		return
	}
	defer c.Close()

	err = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	b := make([]byte, 2)
	i, err := c.Read(b)

	switch {
	case err != nil:
		log.Printf("connection %d read error: %v", n, err)
		return
	case i != len(b):
		log.Printf("connection %d short read: %v", b)
		return
	case string(b) != "Hi":
		log.Printf("connection %d corrupted read: %v", b)
		return
	}

	log.Printf("connection %d: connected: %s -> %s", n, c.LocalAddr(), c.RemoteAddr())

	<- done

	_, err = c.Write([]byte("Hi"))
	if err != nil {
		log.Printf("%s: failed to write: %v", c.RemoteAddr(), err)
		c.Close()
		return
	}
}

func main() {
	c := make(chan struct{})

	for i := uint(0); i < *n; i++ {
		go connect(i, c)
		time.Sleep(100 * time.Millisecond)
	}

	close(c)
	time.Sleep(2 * time.Second)
}
