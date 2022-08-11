package main

import (
	"crypto/tls"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//get_with_new_transport()
	//with_new_transport()
	with_new_transport_close()

}
func with_new_transport() {
	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}

		req, _ := http.NewRequest("GET", "https://google.com", nil)
		client.Do(req)
	}
}

//close idle connection explcitily
func with_new_transport_close() {
	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}

		req, _ := http.NewRequest("GET", "https://google.com", nil)
		client.Do(req)
		client.CloseIdleConnections()
	}
}

//creae new transport for each query, using retryablehttp
func retry_with_new_transport() {
	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		retryClient := retryablehttp.NewClient()
		retryClient.HTTPClient.Transport = tr

		req, _ := retryablehttp.NewRequest("GET", "https://google.com", nil)
		retryClient.Do(req)
	}
}

//reuse transport for each query, using retryablehttp
func retry_reuse_transport() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for {
		retryClient := retryablehttp.NewClient()
		retryClient.HTTPClient.Transport = tr

		req, _ := retryablehttp.NewRequest("GET", "https://google.com", nil)
		retryClient.Do(req)
	}
}
