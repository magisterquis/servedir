package main

/*
 * servedir.go
 * Small program to serve http from a directory
 * By J. Stuart McMurray
 * Created 20140815
 * Last Modified 20160326
 */

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
)

func main() {
	var (
		no4 = flag.Bool(
			"4",
			false,
			"Don't listen on IPv4",
		)
		no6 = flag.Bool(
			"6",
			false,
			"Don't listen on IPv6",
		)
		httpAddr = flag.String(
			"http",
			"8080",
			"HTTP [`address` and] port",
		)
		dir = flag.String(
			"dir",
			".",
			"Files in this `directory` will be served",
		)
		cert = flag.String(
			"cert",
			"cert.pem",
			"HTTPS `certificate` (ignored if -nohttps is given)",
		)
		key = flag.String(
			"key",
			"key.pem",
			"HTTPS `key` (ignored if -nohttps is given)",
		)
		noHTTPS = flag.Bool(
			"nohttps",
			false,
			"Do not handle HTTPS requests",
		)
		httpsAddr = flag.String(
			"https",
			":4433",
			"HTTPS listen [`address` and] port",
		)
		noHTTP = flag.Bool(
			"nohttp",
			false,
			"Don't handle HTTP requests",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage %v [options]

Serves files from the given directory.  Content-type will be automatically
determined.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Make sure we're listening on something */
	if *noHTTP && *noHTTPS {
		log.Fatalf("Both -nohttp and -nohttps can't be given")
	}
	if *no4 && *no6 {
		log.Fatalf("Both -no4 and -no6 can't be given")
	}

	/* Make sure the served directory is a directory */
	if err := verifyDir(*dir); nil != err {
		log.Fatalf("Cannot serve files from %v: %v", *dir, err)
	}
	log.Printf("Serving files from %v", *dir)

	/* Set up a handler for the directory */
	http.Handle(
		"/",
		handlers.CombinedLoggingHandler(
			os.Stderr,
			http.FileServer(http.Dir(*dir)),
		),
	)

	/* Fix up the addresses */
	pa := parsePort(*httpAddr)
	ta := parsePort(*httpsAddr)

	/* WaitGroup to signal end of service */
	wg := &sync.WaitGroup{}

	wg.Add(4)
	go serve("tcp4", pa, !(*noHTTP || *no4), wg)
	go serve("tcp6", pa, !(*noHTTP || *no6), wg)
	go serveTLS("tcp4", ta, !(*noHTTPS || *no4), *cert, *key, wg)
	go serveTLS("tcp6", ta, !(*noHTTPS || *no4), *cert, *key, wg)
	wg.Wait()
	log.Fatalf("Done")
}

/* verifyDir returns nil iff dir is a directory. */
func verifyDir(dir string) error {
	s, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("%v is not a directory", dir)
	}
	return nil
}

/* parseAddress turns a port-only address into a :port address.  If the address
isn't just a number, p is returned */
func parsePort(p string) string {
	if i, err := strconv.Atoi(p); nil == err {
		return fmt.Sprintf(":%v", i)
	}
	return p
}

/* serve serves up HTTP on proto and addr if ok is true.  It calls wg's Done
method before returning. */
func serve(proto, addr string, ok bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if !ok {
		return
	}
	/* Listen on the address */
	l, err := net.Listen(proto, addr)
	if nil != err {
		log.Printf(
			"Unable to listen on %v (%v): %v",
			addr,
			proto,
			err,
		)
		return
	}
	log.Printf("Listening on %v (%v) for HTTP requests", l.Addr(), proto)
	/* Serve */
	log.Printf(
		"Error serving HTTP requests on %v (%v): %v",
		l.Addr(),
		proto,
		http.Serve(l, nil),
	)
}

/* serveTLS serves up HTTPS requests on proto and addr using cert and key if ok
is true.  It calls wg's Done method before returning. */
func serveTLS(
	proto, addr string,
	ok bool,
	cert, key string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	/* Read cert and key */
	c, err := tls.LoadX509KeyPair(cert, key)
	if nil != err {
		log.Printf(
			"Unable to load TLS key from %v and %v for service "+
				"on %v (%v): %v",
			cert, key,
			addr, proto,
			err,
		)
		return
	}

	/* TLS config */
	conf := &tls.Config{
		Certificates: []tls.Certificate{c},
	}

	/* TLS listener */
	l, err := tls.Listen(proto, addr, conf)
	if nil != err {
		log.Printf("Unable to listen on %v (%v) for HTTPS "+
			"requests: %v",
			addr, proto,
			err,
		)
		return
	}
	log.Printf("Listening on %v (%v) for HTTPS requests", l.Addr(), proto)

	/* Serve */
	log.Printf(
		"Error serving HTTP requests on %v (%v): %v",
		l.Addr(),
		proto,
		http.Serve(l, nil),
	)
}
