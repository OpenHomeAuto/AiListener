package webhook

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/OpenHomeAuto/AiListener/pkg/dflow"
	"github.com/OpenHomeAuto/AiListener/pkg/util"
	"github.com/gorilla/mux"
	df "github.com/leboncoin/dialogflow-go-webhook"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"strconv"
	"time"
)

func makeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func Start() {
	var m *autocert.Manager
	var r = mux.NewRouter()

	r.HandleFunc("/webhook", MessagesEndPoint).Methods("POST")
	httpsserver := makeServerFromMux(r)

	dataDir := "."
	m = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(dataDir),
	}

	httpserver := makeServerFromMux(m.HTTPHandler(nil))

	httpserver.Addr = ":80"

	go func() {
		fmt.Printf("Starting HTTP server on %s\n", httpserver.Addr)
		err := httpserver.ListenAndServe()
		if err != nil {
			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		}
	}()

	tlsConfig := &tls.Config{
		Rand:           rand.Reader,
		Time:           time.Now,
		NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
		MinVersion:     tls.VersionTLS12,
		GetCertificate: m.GetCertificate,
	}

	httpsserver.Addr = ":" + strconv.Itoa(*util.HTTPSPort)
	httpsserver.TLSConfig = tlsConfig

	fmt.Printf("Starting HTTPS server on %s\n", httpsserver.Addr)
	err := httpsserver.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
	}
}

func MessagesEndPoint(rw http.ResponseWriter, req *http.Request) {
	var err error
	var dfr *df.Request
	//var p params

	decoder := json.NewDecoder(req.Body)
	if err = decoder.Decode(&dfr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Filter on action, using a switch for example

	// Retrieve the params of the request
	/*if err = dfr.GetParams(&p); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve a specific context
	if err = dfr.GetContext("my-awesome-context", &p); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}*/

	_, err = dflow.DoSignIn(dfr.Session)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Do things with the context you just retrieved
	dff := &df.Fulfillment{
		FulfillmentMessages: df.Messages{
			df.ForGoogle(df.SingleSimpleResponse("hello", "hello")),
			{RichMessage: df.Text{Text: []string{"hello"}}},
		},
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(dff)
}
