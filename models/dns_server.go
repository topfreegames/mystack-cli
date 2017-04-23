package models

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

// DNSServer struct
type DNSServer struct {
	Domains      []string
	forwardToDNS string
	logger       *logrus.Logger
	PointTo      string
	port         int
}

// NewDNSServer ctor
func NewDNSServer(domains []string, forwardToDNS, routerURL string, port int, logger *logrus.Logger) (*DNSServer, error) {
	d := &DNSServer{
		Domains:      domains,
		forwardToDNS: forwardToDNS,
		logger:       logger,
		port:         port,
	}
	err := d.configureDNSServer(routerURL)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DNSServer) configureDNSServer(routerURL string) error {
	uri, err := url.Parse(routerURL)
	if err != nil {
		return err
	}
	host := uri.Hostname()
	routerIPAddr, err := net.LookupHost(host)
	if err != nil {
		return err
	}
	d.PointTo = routerIPAddr[0]
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func (d *DNSServer) parseQuery(m *dns.Msg) {
	l := d.logger.WithFields(logrus.Fields{
		"source": "DNSServer",
	})
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			l.Debugf("Query for %s", q.Name)
			if stringInSlice(q.Name, d.Domains) {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, d.PointTo))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			} else {
				nm := new(dns.Msg)
				nm.RecursionDesired = true
				nm.SetQuestion(q.Name, q.Qtype)
				r, err := dns.Exchange(nm, d.forwardToDNS)
				if err != nil {
					fmt.Print(err.Error())
				}
				m.Answer = r.Answer
			}
		}
	}
}

func (d *DNSServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		d.parseQuery(m)
	}

	w.WriteMsg(m)
}

// Serve starts the local dns server
func (d *DNSServer) Serve() {
	l := d.logger.WithFields(logrus.Fields{
		"source":       "DNSServer",
		"forwardToDNS": d.forwardToDNS,
		"domains":      d.Domains,
		"pointTo":      d.PointTo,
	})
	// attach request handler func
	dns.HandleFunc(".", d.handleDNSRequest)

	// start server on port 53, meaning we will need root permission
	server := &dns.Server{Addr: ":" + strconv.Itoa(d.port), Net: "udp"}
	l.Infof("DNS Server listening at 0.0.0.0:%d\n", d.port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("failed to start server: %s\n ", err.Error())
	}
}
