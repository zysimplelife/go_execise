package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type CertificateWatcher struct {
	needReload chan bool
	cert       string
	key        string
}

func (watcher *CertificateWatcher) watch() {
	watcher.needReload = make(chan bool, 1)
	fswatcher, _ := fsnotify.NewWatcher()
	go func() {
		defer fswatcher.Close()
		for {
			select {
			case event, ok := <-fswatcher.Events:
				if !ok {
					log.Println("error happened")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("event:", event)
					watcher.waitUntilDone(fswatcher)
				}
			case err, ok := <-fswatcher.Errors:
				if !ok {
					log.Println("error:", err)
					return
				}
			}
		}
	}()
	fswatcher.Add(watcher.cert)
	fswatcher.Add(watcher.key)
	log.Println("Started watching")
	return
}

func (watcher CertificateWatcher) waitUntilDone(fswatcher *fsnotify.Watcher) {
	for {
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-fswatcher.Events:
		case <-fswatcher.Errors:
		case <-timer.C:
			log.Println("file change is done.")
			watcher.needReload <- true
			return
		}
		timer.Stop()
	}
}
