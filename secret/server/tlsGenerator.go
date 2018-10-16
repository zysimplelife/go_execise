package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func generateSelfSignCert() (err error) {
	cr := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization: []string{"Yong"},
			Country:      []string{"Sweden"},
			Province:     []string{"Stockholm"},
			Locality:     []string{"Stockholm"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// append ip address
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}

		cr.IPAddresses = append(cr.IPAddresses, ip)
		cr.DNSNames = append(cr.DNSNames, ip.String())
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	// Self sign certificate
	ca_b, err := x509.CreateCertificate(rand.Reader, cr, cr, pub, priv)

	certFile, err := os.Create("server.crt")
	defer certFile.Close()
	log.Println("start generating certificate file")
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})

	// Private key
	keyOut, err := os.OpenFile("server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer keyOut.Close()
	log.Println("start generating key file")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return
}
