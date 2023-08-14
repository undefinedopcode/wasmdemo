package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"path/filepath"

	"log"

	"github.com/gorilla/mux"
)

func fileService(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//log.Println("in serveFunc")

		vars := mux.Vars(r)

		file := "/" + vars["file"]
		if file == "/" {
			file = "/index.html"
		}

		ext := filepath.Ext(file)
		if ext == "" {
			file = "/" + strings.Trim(file, "/") + "/index.html"
		}

		file = baseDir + file

		var data []byte
		var err error

		data, err = ioutil.ReadFile(file)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		case ".otf":
			w.Header().Set("Content-Type", "application/x-font-opentype")
		case ".ttf":
			w.Header().Set("Content-Type", "application/x-font-ttf")
		case ".woff":
			w.Header().Set("Content-Type", "application/x-font-woff")
		case ".wasm":
			w.Header().Set("Content-Type", "application/wasm")
		default:
			w.Header().Set("Content-Type", "text/html")
		}

		w.Write(data)

	}
}

func StartWebServer(hostport string) error {
	r := mux.NewRouter()
	r.Handle("/{file:.*}", fileService("./public")).Methods("GET")
	r.Handle("/", fileService("./public")).Methods("GET")

	var tlsCfg *tls.Config

	if *tlsCert != "" && *tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
		if err != nil {
			log.Fatalf("Failed to load certs: %v", err)
		}
		tlsCfg = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		log.Printf("Starting in secure mode (TLS)")
	} else {
		log.Printf("Starting in basic non-secure mode... use -tls-cert, -tls-key for secure")
	}

	srv := &http.Server{
		Addr: hostport,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 70,
		ReadTimeout:  time.Second * 70,
		IdleTimeout:  time.Second * 120,
		Handler:      r, // Pass our instance of gorilla/mux in.
		TLSConfig:    tlsCfg,
	}

	log.Printf("Starting web server on %s", hostport)

	if srv.TLSConfig != nil {
		return srv.ListenAndServeTLS("", "")
	}
	return srv.ListenAndServe()
}

var port = flag.String("hostport", ":8080", "Host port to run on")
var tlsCert = flag.String("tls-cert", "", "TLS certificate (requires -tls-key)")
var tlsKey = flag.String("tls-key", "", "TLS key file (requires -tls-cert)")

func main() {
	flag.Parse()
	log.Fatal(StartWebServer(*port))
}
