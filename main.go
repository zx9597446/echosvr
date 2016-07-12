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

var sem = make(chan int, 100)

func main() {
	flag.Parse()
	if *protocol == "tcp" {
		tcpServer()
	} else if *protocol == "udp" {
		udpServer()
	}
}

func udpServer() {
	udpaddr, err := net.ResolveUDPAddr(*protocol, *addr)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("listening on ", udpaddr)
	udpconn, err := net.ListenUDP(*protocol, udpaddr)
	if err != nil {
		log.Panicln(err)
	}
	defer udpconn.Close()
	for {
		sem <- 1
		go udpEcho(udpconn)
	}
}

func udpEcho(con net.PacketConn) {
	defer func() { <-sem }()
	buf := make([]byte, 4096)
	nr, addr, err := con.ReadFrom(buf)
	if err != nil {
		log.Print(err)
		return
	}
	nw, err := con.WriteTo(buf[:nr], addr)
	if err != nil {
		log.Print(err)
		return
	}
	if nw != nr {
		log.Printf("received %d bytes but sent %d\n", nr, nw)
	}
}

func tcpServer() {
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
