package http

import (
	"context"
	"fmt"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/rest"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

// Server interface for starting and stopping the HTTP server
type Server interface {
	Start()
	Stop()
}

// GostServer is the type that contains all of the relevant information to set
// up the GOST HTTP Server
type GostServer struct {
	host       string      // Hostname for example "localhost" or "192.168.1.14"
	port       int         // Portnumber where you want to run your http server on
	api        *models.API // Sensorthings api to interact with from the HttpServer
	https      bool
	httpsCert  string
	httpsKey   string
	httpServer *http.Server
}

// CreateServer initialises a new GOST HTTPServer based on the given parameters
func CreateServer(host string, port int, api *models.API, https bool, httpsCert, httpsKey string) Server {
	router := CreateRouter(api)
	return &GostServer{
		host:      host,
		port:      port,
		api:       api,
		https:     https,
		httpsCert: httpsCert,
		httpsKey:  httpsKey,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, strconv.Itoa(port)),
			Handler:      PostProcessHandler(LowerCaseURI(router)),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
}

// Start command to start the GOST HTTPServer
func (s *GostServer) Start() {
	t := "HTTP"
	if s.https {
		t = "HTTPS"
	}
	log.Printf("Started GOST %v Server on %v:%v", t, s.host, s.port)

	var err error
	if s.https {
		err = s.httpServer.ListenAndServeTLS(s.httpsCert, s.httpsKey)
	} else {
		err = s.httpServer.ListenAndServe()
	}

	if err != nil {
		log.Fatal(err)
		return
	}
}

// Stop command to stop the GOST HTTP server, currently not supported
func (s *GostServer) Stop() {
	if s.httpServer != nil {
		log.Print("Stopping HTTP(S) Server")
		s.httpServer.Shutdown(context.Background())
	}
}

// LowerCaseURI is a middleware function that lower cases the url path
func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		lowerCasePath := strings.ToLower(r.URL.Path)
		// temporarily disabled checking on paths due to problems in serving with /$value
		/*
			api := *s.api
			split := strings.Split(lowerCasePath, "/")
			for i, s := range split {
				if len(s) == 0 || i+1 == len(split) {
					continue
				}

				found := false
				for _, a := range api.GetAcceptedPaths() {
					if strings.HasPrefix(s, a) {
						found = true
					}
				}

				if !found {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(""))
					return
				}
			}
		*/

		r.URL.RawPath = r.URL.Path
		r.URL.Path = lowerCasePath
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// PostProcessHandler runs after all the other handlers and can be used to modify
// the response. In this case modify links (due to proxy running) or handle CORS functionality
// Basically we catch all the response using httptest.NewRecorder (headers + body), modify and
// write to response. Is this a right approach?
func PostProcessHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		orig_uri := rest.ExternalURI
		forwarded_uri := r.Header.Get("X-Forwarded-For")

		rec := httptest.NewRecorder()

		// first run the next handler and get results
		h.ServeHTTP(rec, r)

		// read response body and replace links
		bytes := rec.Body.Bytes()
		s := string(bytes)
		if len(s) > 0 {
			if len(forwarded_uri) > 0 {
				s = strings.Replace(s, orig_uri, forwarded_uri, -1)
			}

		}
		// handle headers too...
		for k, v := range rec.HeaderMap {
			val := v[0]

			// if there is a location header and a proxy running, change
			// the location url too
			if k == "Location" && len(forwarded_uri) > 0 {
				val = strings.Replace(val, orig_uri, forwarded_uri, -1)
			}

			// add the header to response
			w.Header().Add(k, val)
		}

		// handle status code
		w.WriteHeader(rec.Code)

		// now add CORS response header
		w.Header().Add("Access-Control-Allow-Origin", "*")

		// write modified response
		w.Write([]byte(s))
	}
	return http.HandlerFunc(fn)
}
