package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	ip4flag := flag.Bool("ip4", false, "Use IPv4 (default).")
	ip6flag := flag.Bool("ip6", false, "Use IPv6.")
	addr := flag.String("addr", ":8080", "[IP Address and] port on "+
		"which to listen.")
	dir := flag.String("dir", "./", "Directory to serve.")
	flag.Parse()

	/* Make sure argv[1] is a directory */
	s, err := os.Stat(*dir)
	if err != nil {
		log.Printf("Unable to stat %v: %v", *dir, err)
		os.Exit(-3)
	}
	if !s.IsDir() {
		log.Printf("%v is not a directory", *dir)
		os.Exit(-2)
	}
	/* Work out the port */
	if _, err := strconv.Atoi(*addr); nil == err {
		*addr = fmt.Sprintf(":%v", *addr)
	}
	/* Work out the network on which to listen */
	network := "tcp4"
	if *ip4flag && *ip6flag {
		network = "tcp"
	} else if *ip6flag {
		network = "tcp6"
	}
	/* Make a listener */
	ln, err := net.Listen(network, *addr)
	if err != nil {
		log.Printf("Unable to listen: %v", err)
		os.Exit(-1)
	}
	/* Make a server */
	server := &http.Server{Addr: *addr,
		Handler: http.FileServer(http.Dir(*dir))}
	log.Printf("Serving up %v on %v\n", *dir, ln.Addr())
	/* Start serving */
	log.Fatal(server.Serve(ln))
}
