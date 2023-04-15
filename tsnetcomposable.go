package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"tailscale.com/tsnet"
)

var (
	enable_funnel_string = os.Getenv("TSNET_ENABLE_FUNNEL")
	enable_funnel        = false
	redirect_to          = os.Getenv("TSNET_HTTP_REDIRECT_URL")
	addr                 = os.Getenv("TSNET_LISTEN_ADDR")
	custom_hostname      = os.Getenv("TSNET_CUSTOM_HOSTNAME")
	proxy_to             = os.Getenv("TSNET_PROXY_TO_URL")
)

func main() {
	flag.Parse()

	if strings.ToLower(enable_funnel_string) == "true" {
		enable_funnel = true
		log.Println("running as a PUBLIC funneled node.")
	}

	if addr == "" {
		log.Println("no listen address provided. assuming :443")
		addr = ":443"
	}
	if custom_hostname == "" {
		log.Fatal("missing hostname")
	}

	//TODO can theoretically just get this URL instead of having it provided. but :shrug:
	//this is easier for now
	if proxy_to == "" {
		log.Fatal("missing proxy-url. likely https://container")
	}
	proxy_to_url, err := url.Parse(proxy_to)
	if err != nil {
		log.Fatalln("failed to parse proxy url")
	}

	log.Println("starting")

	s := &tsnet.Server{
		Hostname: custom_hostname,
	}

	defer s.Close()

	var ln net.Listener
	if enable_funnel {
		ln, err = s.ListenFunnel("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
	} else if addr == ":443" {
		ln, err = s.ListenTLS("tcp", ":443")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ln, err = s.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer ln.Close()

	// lc, err := s.LocalClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if addr == ":443" {
	// 	ln = tls.NewListener(ln, &tls.Config{
	// 		GetCertificate: lc.GetCertificate,
	// 	})
	// }

	proxy := httputil.NewSingleHostReverseProxy(proxy_to_url)

	log.Fatal(http.Serve(ln, proxy))
}
