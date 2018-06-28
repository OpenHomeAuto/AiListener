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

type OriginalResp struct {
	Source  string `json:"source"`
	Version string `json:"version"`
	Payload struct {
		IsInSandbox bool `json:"isInSandbox"`
		Surface     struct {
			Capabilities []struct {
				Name string `json:"name"`
			} `json:"capabilities"`
		} `json:"surface"`
		RequestType string `json:"requestType"`
		Inputs      []struct {
			RawInputs []struct {
				Query     string `json:"query"`
				InputType string `json:"inputType"`
			} `json:"rawInputs"`
			Arguments []struct {
				RawText   string `json:"rawText"`
				TextValue string `json:"textValue"`
				Name      string `json:"name"`
			} `json:"arguments"`
			Intent string `json:"intent"`
		} `json:"inputs"`
		User struct {
			LastSeen    time.Time `json:"lastSeen"`
			Locale      string    `json:"locale"`
			UserID      string    `json:"userId"`
			AccessToken string    `json:"accessToken"`
		} `json:"user"`
		Conversation struct {
			ConversationID    string `json:"conversationId"`
			Type              string `json:"type"`
			ConversationToken string `json:"conversationToken"`
		} `json:"conversation"`
		AvailableSurfaces []struct {
			Capabilities []struct {
				Name string `json:"name"`
			} `json:"capabilities"`
		} `json:"availableSurfaces"`
	} `json:"payload"`
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
	switch dfr.QueryResult.Action {
	case "play_music":
		log.Println(dfr.QueryResult)

		var oresp *OriginalResp
		if err = json.Unmarshal(dfr.OriginalDetectIntentRequest, &oresp); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if oresp.Payload.User.AccessToken == "" {
			resp := dflow.DoSignIn()
			// Do things with the context you just retrieved
			dff := &df.Fulfillment{
				FollowupEventInput: resp,
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(dff)

			return
		}

		// Do things with the context you just retrieved
		dff := &df.Fulfillment{
			FulfillmentMessages: df.Messages{
				df.ForGoogle(df.SingleSimpleResponse("Starting Music", "Starting Music")),
				{RichMessage: df.Text{Text: []string{"Starting Music"}}},
			},
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(dff)
		return
	case "auth":
		var oresp *OriginalResp
		if err = json.Unmarshal(dfr.OriginalDetectIntentRequest, &oresp); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("OrespJ: \n %+v\n", oresp)
		respPJ, _ := dfr.QueryResult.Parameters.MarshalJSON()
		log.Println("respPJ: ", string(respPJ))
		for _, v := range dfr.QueryResult.OutputContexts {
			log.Println(v)
			respOCPJ, _ := v.Parameters.MarshalJSON()
			log.Println("respOCPJ: ", string(respOCPJ))
		}
		/*
			type params struct {

			}

			var p *params

			if err = dfr.GetParams(&p); err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Println(p)
		*/

		dff := &df.Fulfillment{
			FulfillmentMessages: df.Messages{
				df.ForGoogle(df.SingleSimpleResponse("Authentication finished", "Authentication finished")),
				{RichMessage: df.Text{Text: []string{"Authentication finished"}}},
			},
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(dff)
		return
	default:
		log.Println(dfr.QueryResult)

		dff := &df.Fulfillment{
			FulfillmentMessages: df.Messages{
				df.ForGoogle(df.SingleSimpleResponse("default", "default")),
				{RichMessage: df.Text{Text: []string{"default"}}},
			},
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(dff)
		return
	}

}
