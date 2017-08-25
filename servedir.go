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
	"sync"

	"github.com/gorilla/handlers"
)

func main() {
	var (
		httpAddr = flag.String(
			"http",
			"0.0.0.0:8080",
			"HTTP listen `address`",
		)
		dir = flag.String(
			"dir",
			".",
			"Files in this `directory` will be served",
		)
		cert = flag.String(
			"cert",
			"cert.pem",
			"HTTPS certificate `file`",
		)
		key = flag.String(
			"key",
			"key.pem",
			"HTTPS key `file`",
		)
		httpsAddr = flag.String(
			"https",
			"0.0.0.0:4433",
			"HTTPS listen `address`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage %v [options]

Serves files from the given directory.  Content-type will be automatically
determined.

If the address given to -http or -https is no http or https will not be
served, respectively.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Make sure we have at least one address */
	if "" == *httpAddr && "" == *httpsAddr {
		log.Fatalf("No listen addresses specified")
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

	/* Listen and serve */
	wg := &sync.WaitGroup{}
	if "no" != *httpAddr {
		wg.Add(1)
		go serveHTTP(*httpAddr, wg)
	}
	if "no" != *httpsAddr {
		wg.Add(1)
		go serveHTTPS(*httpsAddr, *cert, *key, wg)
	}
	wg.Wait()
	log.Printf("Done.")
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

/* serveHTTP serves files over HTTP */
func serveHTTP(addr string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Listen */
	l, err := net.Listen("tcp", addr)
	if nil != err {
		log.Fatalf(
			"Unable to listen for HTTP requests to %q: %v",
			addr,
			err,
		)
	}
	log.Printf("Listening for HTTP connections to %v", l.Addr())
	/* Serve */
	if err := http.Serve(l, nil); nil != err {
		log.Fatalf(
			"Error serving HTTP requests to %v: %v",
			l.Addr(),
			err,
		)
	}
}

/* serveHTTPS serves files over HTTPS */
func serveHTTPS(addr, cert, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Make TLS */
	c, err := tls.LoadX509KeyPair(cert, key)
	if nil != err {
		log.Fatalf(
			"Unable to load keypair from %q and %q: %v",
			cert,
			key,
			err,
		)
	}
	/* Listen */
	l, err := tls.Listen("tcp", addr, &tls.Config{
		Certificates: []tls.Certificate{c},
	})
	if nil != err {
		log.Fatalf(
			"Unable to listen for HTTPS connections to %q: %v",
			addr,
			err,
		)
	}
	log.Printf("Listening for HTTPS connections to %v", l.Addr())
	if err := http.Serve(l, nil); nil != err {
		log.Fatalf(
			"Error serving HTTPS requests to %v: %v",
			l.Addr(),
			err,
		)
	}
}
