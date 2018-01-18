package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/pelletier/go-toml"
)

type Domain struct {
	Name string `toml:"Name"`
	Path string `toml:"Path"`
	Cert string `toml:"Certificate"`
	Key  string `toml:"PrivateKey"`
}

type Config struct {
	Port      string   `toml:"PlainPort"`
	SSLPort   string   `toml:"SSLPort"`
	Log       string   `toml:"Log"`
	Domains   []Domain `toml:"Domains"`
	tlsConfig *tls.Config
}

type Request struct {
	handler http.Handler
	path    string
}

func ReadConfig() {
	c := &Config{}

	// reading config file
	data, err := ioutil.ReadFile(confile)
	if err != nil {
		ch <- err.Error()
		return
	}

	// decoding configuration
	err = toml.Unmarshal(data, c)
	if err != nil {
		ch <- err.Error()
		return
	}

	// creating tls config for https
	c.tlsConfig = &tls.Config{
		PreferServerCipherSuites: true,
		Certificates:             make([]tls.Certificate, 0),
	}

	// configuring domains
	for _, domain := range c.Domains {
		domains[domain.Name] = http.FileServer(
			http.Dir(domain.Path),
		)
		if useTLS {
			// configuring certs for https
			cert, err := tls.LoadX509KeyPair(
				domain.Cert, domain.Key,
			)
			if err != nil {
				ch <- err.Error()
			} else {
				c.tlsConfig.Certificates = append(c.tlsConfig.Certificates, cert)
			}
		}
	}
	if useTLS {
		c.tlsConfig.BuildNameToCertificate()
	}
	config = c
}
