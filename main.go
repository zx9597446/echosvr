package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var addr = flag.String("addr", ":11111", "addr to listen")

func main() {
	server, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("listening on ", *addr)
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				log.Println(err)
			}
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	io.Copy(io.MultiWriter(client, os.Stdout), client)
}
