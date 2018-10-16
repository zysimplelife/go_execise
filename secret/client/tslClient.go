package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
)

var (
	caFile = flag.String("CA", "../server/server.crt", "A PEM eoncoded CA's certificate file.")
)

func main() {
	log.SetFlags(log.Lshortfile)
	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	conf := &tls.Config{RootCAs: caCertPool}

	conn, err := tls.Dial("tcp", "127.0.0.1:7443", conf)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	n, err := conn.Write([]byte("Hello \n"))
	if err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}
