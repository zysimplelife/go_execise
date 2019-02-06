package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Server is starting")

	// start ticker to update certifcate
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			generateSelfSignCert()
		}
	}()

	//server
	server := ReloadableListenerV2{
		cert: "./server.crt",
		key:  "./server.key",
	}

	// watching certificate change
	watcher := CertificateWatcher{
		cert: "./server.crt",
		key:  "./server.key",
	}
	watcher.watch()
	generateSelfSignCert()
	ctx, cancel := context.WithCancel(context.Background())
	server.serve(ctx)
	for {
		select {
		case <-watcher.needReload:
			cancel()
			<-server.done
			server.serve(ctx)
		}
	}

	time.Sleep(60 * time.Second)
}
