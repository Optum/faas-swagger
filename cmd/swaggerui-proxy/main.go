package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	_ "github.com/optum/faas-swagger/static/statik"
	"github.com/rakyll/statik/fs"
)

func main() {

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", http.FileServer(statikFS)))

	//reverse proxy openfaas request
	http.HandleFunc("/", serveReverseProxy)

	http.ListenAndServe(":8080", nil)

}

func serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	gateway, _ := os.LookupEnv("OPENFAAS_GATEWAY")
	u, _ := url.Parse(gateway)
	u.Path = "/function"

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = u.Host

	proxy.ServeHTTP(res, req)
}
