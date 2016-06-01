package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var addr = flag.String("a", ":11111", "addr to listen")
var protocol = flag.String("p", "tcp", "protocol, default tcp")
var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()
	server, err := net.Listen(*protocol, *addr)
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
	if *verbose {
		io.Copy(io.MultiWriter(client, os.Stdout), client)
	} else {
		io.Copy(client, client)
	}
}
