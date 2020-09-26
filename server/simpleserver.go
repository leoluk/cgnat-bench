package main

import (
	"flag"
	"log"
	"net"
	"time"
)

var (
	addr = flag.String("addr", "[::]:8080", "Listen address")
)

func init() {
	flag.Parse()
}

func main() {
	s, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listening on %v", s.Addr())

	for {
		c, err := s.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			return
		}
		defer c.Close()

		go func(c net.Conn) {
			log.Printf("got connection from %s", c.RemoteAddr())

			err = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Printf("%s: failed to set write deadline: %v", c.RemoteAddr(), err)
				c.Close()
				return
			}

			_, err := c.Write([]byte("Hi"))
			if err != nil {
				log.Printf("%s: failed to write: %v", c.RemoteAddr(), err)
				c.Close()
				return
			}

			err = c.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				log.Printf("%s: failed to set read deadline: %v", c.RemoteAddr(), err)
				c.Close()
				return
			}

			b := make([]byte, 2)
			_, err = c.Read(b)
			if err != nil {
				log.Printf("%s: read failed: %v", c.RemoteAddr(), err)
				return
			}

			//log.Printf("%s: read %d bytes: %v", c.RemoteAddr(), i, string(b))
			log.Printf("%s: closing connection", c.RemoteAddr())
		}(c)
	}
}
