package main

import (
	"crypto/tls"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{}

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
