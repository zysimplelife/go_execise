package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fsnotify/fsnotify"
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
			go handleConnection(conn)
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

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Server is starting")

	// start ticker to update certifcate
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			generateSelfSignCert()
		}
	}()

	// start server
	server := ReloadableListener{
		cert: "./server.crt",
		key:  "./server.key",
	}

	// watching certificate change
	log.Println("Start receiving message from watcher")
	needReload := watchCertficate()

	for {
		select {
		case <-needReload:
			server.restart()
		}
	}

	time.Sleep(60 * time.Second)
}

func watchCertficate() (needReload chan bool) {
	needReload = make(chan bool, 1)
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("error happened")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("event:", event)
					waitUntilDone(watcher, needReload)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("error:", err)
					return
				}
			}
		}
	}()
	watcher.Add("server.crt")
	watcher.Add("server.key")
	return
}

func waitUntilDone(watcher *fsnotify.Watcher, needReload chan bool) {
	for {
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-watcher.Events:
		case <-watcher.Errors:
		case <-timer.C:
			log.Println("file change is done.")
			needReload <- true
			return
		}
		timer.Stop()
	}
}

func handleConnection(conn net.Conn) {
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
