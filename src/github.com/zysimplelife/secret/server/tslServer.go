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
	ln, _ := createListener()
	done := make(chan bool, 1)
	go serve(done, ln)

	// watching certificate change
	log.Println("Start receiving message from watcher")
	needReload := watchCertficate()

	for {
		select {
		case <-needReload:
			log.Println("Restarting server....")
			stopServing(done, ln)
			ln, _ = createListener()
			go serve(done, ln)
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
			log.Println("no new event update, file change is done.")
			needReload <- true
			return
		}
		timer.Stop()
	}
}

func createListener() (ln net.Listener, err error) {
	cer, err := tls.LoadX509KeyPair("./server.crt", "./server.key")
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

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err = tls.Listen("tcp", ":7443", config)
	if err != nil {
		log.Println("some error happened ", err)
	}
	return
}

func serve(done chan bool, ln net.Listener) {
	//serve incoming request
	log.Println("Serving started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Some error happened", err)
			break
		}
		go handleConnection(conn)
	}
	log.Println("stop serving,release listener ", ln)
	done <- true
}

func stopServing(done chan bool, ln net.Listener) {
	if ln == nil {
		log.Println("listener is nil, no need release")
	}
	ln.Close()
	<-done
	log.Println("Serving stoped")
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
