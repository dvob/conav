package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dvob/conav"
	"github.com/dvob/conav/report/so"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	var (
		tls  bool
		host string
		addr string
	)

	flag.BoolVar(&tls, "tls", false, "use acme to configure certificates and serve tls")
	flag.StringVar(&host, "host", "", "hostname for the certificate")
	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.Parse()

	solothurnCaseReporter := so.NewCaseReporter()
	caseReportHandler := conav.NewCaseReportHanlder(solothurnCaseReporter)
	http.Handle("/", caseReportHandler)

	s := &http.Server{
		Addr: addr,
	}

	if tls {
		if host == "" {
			log.Fatal("flag -host required with -tls")
		}
		m := &autocert.Manager{
			Cache:      autocert.DirCache("acme"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(host),
		}
		s.TLSConfig = m.TLSConfig()
		log.Fatal(s.ListenAndServeTLS("", ""))
		return
	}

	log.Fatal(s.ListenAndServe())

}
