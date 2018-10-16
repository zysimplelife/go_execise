package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
)

type ReloadableListener struct {
	done chan bool
	cert string
	key  string
	ln   net.Listener
}

func (rl ReloadableListener) loadCertificate() (tlsConfig *tls.Config, err error) {
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

func (rl *ReloadableListener) serve() {
	log.Println("Starting server")
	config, err := rl.loadCertificate()
	if err != nil {
		log.Println("some error happened ", err)
	}
	//serve incoming request
	rl.ln, err = tls.Listen("tcp", ":7443", config)
	if err != nil {
		log.Println("some error happened ", err)
	}
	rl.done = make(chan bool, 1)
	go func() {
		for {
			log.Println("Start accepting request with addr", rl.ln.Addr)
			conn, err := rl.ln.Accept()
			if err != nil {
				log.Println("Some error happened", err)
				break
			}
			go rl.handleConnection(conn)
		}
		log.Println("stop serving,release listener ", rl.ln)
		rl.done <- true
	}()
	log.Println("Serving started")
}

func (rl *ReloadableListener) stopServing() {
	if rl.ln == nil {
		log.Println("listener is nil, no need release")
		return
	}
	rl.ln.Close()
	<-rl.done
	rl.ln = nil
	log.Println("Serving stoped")
}

func (rl *ReloadableListener) restart() {
	log.Println("Restarting server....")
	rl.stopServing()
	rl.serve()
}

func (rl ReloadableListener) handleConnection(conn net.Conn) {
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
