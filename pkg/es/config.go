package es

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/olivere/elastic/v7"
)

type Config struct {
	Servers  []string
	Username string
	Password string
	Sniffer  bool // https://github.com/olivere/elastic/wiki/Sniffing
	TLS      TLSConfig
	Version  uint
}

type TLSConfig struct {
	Enabled        bool
	SkipHostVerify bool
	CertPath       string
	KeyPath        string
	CaPath         string
}

func (cfg *Config) getConfigOptions() ([]elastic.ClientOptionFunc, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Servers...),
		elastic.SetSniff(cfg.Sniffer),
		elastic.SetHealthcheck(false),
	}
	httpClient := &http.Client{Timeout: timeout}
	options = append(options, elastic.SetHttpClient(httpClient))
	if cfg.TLS.Enabled {
		ctlsConfig, err := cfg.TLS.createTLSConfig()
		if err != nil {
			return nil, err
		}
		httpClient.Transport = &http.Transport{TLSClientConfig: ctlsConfig}
		return options, nil
	}

	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.TLS.SkipHostVerify,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS12,
		},
	}
	if cfg.TLS.CaPath != "" {
		ctls := &TLSConfig{CaPath: cfg.TLS.CaPath}
		ca, err := ctls.loadCertificate()
		if err != nil {
			return nil, err
		}
		httpTransport.TLSClientConfig.RootCAs = ca
	}

	httpClient.Transport = httpTransport
	if cfg.Username != "" && cfg.Password != "" {
		options = append(options, elastic.SetBasicAuth(cfg.Username, cfg.Password))
	}
	return options, nil
}

func (tlsConfig *TLSConfig) createTLSConfig() (*tls.Config, error) {
	rootCerts, err := tlsConfig.loadCertificate()
	if err != nil {
		return nil, err
	}
	clientPrivateKey, err := tlsConfig.loadPrivateKey()
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            rootCerts,
		Certificates:       []tls.Certificate{*clientPrivateKey},
		InsecureSkipVerify: tlsConfig.SkipHostVerify,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		MinVersion: tls.VersionTLS11,
		MaxVersion: tls.VersionTLS12,
	}, nil

}

func (tlsConfig *TLSConfig) loadCertificate() (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(tlsConfig.CaPath)
	if err != nil {
		return nil, err
	}
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM(caCert)
	return certificates, nil
}

func (tlsConfig *TLSConfig) loadPrivateKey() (*tls.Certificate, error) {
	privateKey, err := tls.LoadX509KeyPair(tlsConfig.CertPath, tlsConfig.KeyPath)
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}
