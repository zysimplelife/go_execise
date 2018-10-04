package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cer, err := tls.LoadX509KeyPair("./server.crt", "./server.key")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Start Server with certificate")
	log.Println("size in the certificate chain =", len(cer.Certificate))

	cerX509, err := x509.ParseCertificate(cer.Certificate[0])
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("certificate is for =", cerX509.Subject)
	log.Println("certificate issuer is  =", cerX509.Issuer)

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":7443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
}
