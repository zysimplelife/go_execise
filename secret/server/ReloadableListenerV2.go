package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
)

type ReloadableListenerV2 struct {
	done chan bool
	cert string
	key  string
}

func (rl ReloadableListenerV2) loadCertificate() (tlsConfig *tls.Config, err error) {
	cer, err := tls.LoadX509KeyPair(rl.cert, rl.key)
	if err != nil {
		log.Println("some error happened ", err)
	}

	log.Println("size in the certificate chain =", len(cer.Certificate))

	cerX509, err := x509.ParseCertificate(cer.Certificate[0])
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("certificate is for ", cerX509.Subject)
	log.Println("certificate issuer is  ", cerX509.Issuer)
	log.Println("IP Addresse is ", cerX509.IPAddresses)

	tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	return
}

func (rl *ReloadableListenerV2) serve(ctx context.Context) {
	log.Println("Starting server")
	config, err := rl.loadCertificate()
	if err != nil {
		log.Println("some error happened ", err)
	}
	//serve incoming request
	ln, err := tls.Listen("tcp", ":7443", config)
	if err != nil {
		log.Println("some error happened ", err)
	}
	rl.done = make(chan bool, 1)

	// start listing
	go func() {
		for {
			log.Println("Start accepting request with addr", ln.Addr)
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Some error happened", err)
				break
			}
			go rl.handleConnection(conn)
		}
		log.Println("stop serving,release listener ", ln)
		rl.done <- true
	}()

	// use context to close server
	go func() {
		for {
			<-ctx.Done()
			log.Println("Server should be done")
			ln.Close()
		}
	}()

	log.Println("Serving started")
}

func (rl ReloadableListenerV2) handleConnection(conn net.Conn) {
	log.Println("new connection is created from =", conn.RemoteAddr)
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
