package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"tailscale.com/tsnet"
)

var (
	enable_redirect_string = os.Getenv("TSNET_ENABLE_REDIRECT")
	enable_redirect        = false
	redirect_to            = os.Getenv("TSNET_HTTP_REDIRECT_URL")
	addr                   = os.Getenv("TSNET_LISTEN_ADDR")
	custom_hostname        = os.Getenv("TSNET_CUSTOM_HOSTNAME")
	proxy_to               = os.Getenv("TSNET_PROXY_TO_URL")
)

func main() {
	flag.Parse()

	// mux := http.NewServeMux()
	if strings.ToLower(enable_redirect_string) == "true" {
		enable_redirect = true
		log.Println("enabling http -> https redirect")
		if redirect_to == "" {
			log.Fatal("missing redirect url. likely <hostname>.something.ts.net")

		}
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
	ln, err := s.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	if enable_redirect {
		go func() {
			ln_redirect, err := s.Listen("tcp", ":80")
			if err != nil {
				log.Fatal(err)
			}
			defer ln_redirect.Close()
			log.Printf("setting up redirect to %s\n", redirect_to)
			log.Fatal(http.Serve(ln_redirect, http.RedirectHandler(redirect_to, http.StatusMovedPermanently)))
		}()
	}

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	if addr == ":443" {
		ln = tls.NewListener(ln, &tls.Config{
			GetCertificate: lc.GetCertificate,
		})
	}

	// log.Fatal(http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "<html><body><h1>Hello, world!</h1>\n")
	// 	fmt.Fprintf(w, "<p>You are <b>%s</b> from <b>%s</b> (%s) (%s)</p>",
	// 		html.EscapeString(who.UserProfile.LoginName),
	// 		html.EscapeString(who.Node.ComputedName),
	// 		r.RemoteAddr,
	// 		html.EscapeString(who.Node.Name))
	// })))
	proxy := httputil.NewSingleHostReverseProxy(proxy_to_url)

	log.Fatal(http.Serve(ln, proxy))
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}
