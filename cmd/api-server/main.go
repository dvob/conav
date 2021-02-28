package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dvob/conav"
	"github.com/dvob/conav/report/so"
	"github.com/gorilla/handlers"
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

	s := &http.Server{
		Addr:    addr,
		Handler: handlers.LoggingHandler(os.Stdout, caseReportHandler),
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
