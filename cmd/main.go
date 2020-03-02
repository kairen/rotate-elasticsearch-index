package main

import (
	goflag "flag"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kairen/rotate-elasticsearch-index/pkg/es"
	flag "github.com/spf13/pflag"
)

var (
	endpoints          []string
	days               int
	indexRegexPatterns []string
	dateFormat         string
	retryCount         int
	basicUsername      string
	basicPassword      string
	sniffer            bool
	tlsEnable          bool
	caPath             string
	certPath           string
	keyPath            string
	skipHostVerify     bool
)

func parseFlags() {
	flag.StringSliceVarP(&endpoints, "endpoints", "", []string{"http://elasticsearch:9200"}, "Endpoints of elasticsearch.")
	flag.StringVarP(&basicUsername, "username", "", os.Getenv("ELASTIC_USERNAME"), "Username for basic auth.")
	flag.StringVarP(&basicPassword, "password", "", os.Getenv("ELASTIC_PASSWORD"), "Password for basic auth.")
	flag.BoolVarP(&sniffer, "sniffer", "", false, "Enable client to use a sniffing process for finding all nodes of your cluster.")
	flag.BoolVarP(&tlsEnable, "tls", "", false, "Enable that servers are TLS.")
	flag.StringVarP(&caPath, "tls.ca", "", "", "SSL Certificate Authority file used to secure elasticsearch communication.")
	flag.StringVarP(&certPath, "tls.cert", "", "", "SSL certification file used to secure elasticsearch communication.")
	flag.StringVarP(&keyPath, "tls.key", "", "", "SSL key file used to secure elasticsearch communication.")
	flag.BoolVarP(&skipHostVerify, "tls.skip-host-verify", "", false, "(insecure) Skip server's certificate chain and host name verification")
	flag.IntVarP(&days, "days", "d", 90, "Days to keep.")
	flag.StringSliceVarP(&indexRegexPatterns, "index-regex-patterns", "", nil, "Index's regex pattern.")
	flag.StringVarP(&dateFormat, "date-format", "", "2006.1.2", "Format template for parsing date.")
	flag.IntVarP(&retryCount, "retry-count", "", 5, "The number of retry for deleting request.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}

func main() {
	defer glog.Flush()
	parseFlags()

	if len(indexRegexPatterns) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := es.Config{
		Servers:  endpoints,
		Username: basicUsername,
		Password: basicPassword,
		Sniffer:  sniffer,
		TLS: es.TLSConfig{
			Enabled:        tlsEnable,
			CaPath:         caPath,
			CertPath:       certPath,
			KeyPath:        keyPath,
			SkipHostVerify: skipHostVerify,
		},
	}
	client, err := es.NewClient(cfg)
	if err != nil {
		glog.Fatalln(err)
	}

	for _, name := range indexRegexPatterns {
		err := retry(time.Second*2, retryCount, func() error {
			return client.RotateIndex(name, days, dateFormat)
		})
		if err != nil {
			glog.Errorf("Failed to remove \"%s\" indices: %s.", name, err)
		}
	}
}

func retry(d time.Duration, attempts int, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}
		glog.Errorf("Error: %s, Retrying in %s. %d Retries remaining.", err, d, attempts-i)
		time.Sleep(d)
	}
	return err
}
