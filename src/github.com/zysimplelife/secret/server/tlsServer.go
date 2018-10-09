package main

import (
	"fmt"
	"log"
	"time"
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

	//server
	server := ReloadableListener{
		cert: "./server.crt",
		key:  "./server.key",
	}

	// watching certificate change
	watcher := CertificateWatcher{
		cert: "./server.crt",
		key:  "./server.key",
	}
	watcher.watch()
	for {
		select {
		case <-watcher.needReload:
			server.restart()
		}
	}

	time.Sleep(60 * time.Second)
}
