package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	gostLog "github.com/gost/server/log"
	"github.com/gost/server/sensorthings/models"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func setupLogger() {
	l, err := gostLog.GetLoggerInstance()
	if err != nil {
		log.Error(err)
	}

	//Setting default fields for main logger
	logger = l.WithFields(log.Fields{"package": "gost.server.http"})
}

// Server interface for starting and stopping the HTTP server
type Server interface {
	Start()
	Stop()
}

// GostServer is the type that contains all of the relevant information to set
// up the GOST HTTP Server
type GostServer struct {
	host       string      // Hostname for example "localhost" or "192.168.1.14"
	port       int         // Port number where you want to run your http server on
	api        *models.API // SensorThings api to interact with from the HttpServer
	https      bool
	httpsCert  string
	httpsKey   string
	httpServer *http.Server
}

// CreateServer initialises a new GOST HTTPServer based on the given parameters
func CreateServer(host string, port int, api *models.API, https bool, httpsCert, httpsKey string) Server {
	setupLogger()
	a := *api
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
			Handler:      PostProcessHandler(LowerCaseURI(router), a.GetConfig().Server.ExternalURI),
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

	logger.Infof("Started GOST %v Server on %v:%v", t, s.host, s.port)

	var err error
	if s.https {
		err = s.httpServer.ListenAndServeTLS(s.httpsCert, s.httpsKey)
	} else {
		err = s.httpServer.ListenAndServe()
	}

	if err != nil {
		logger.Panicf("GOST server not properly stopped: %v", err)
	}
}

// Stop command to stop the GOST HTTP server
func (s *GostServer) Stop() {
	if s.httpServer != nil {
		logger.Info("Stopping HTTP(S) Server")
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
func PostProcessHandler(h http.Handler, externalURI string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if logger.Logger.Level == log.DebugLevel {
			logger.Debugf("%s start: %s", r.Method, r.URL.Path)
			defer gostLog.DebugfWithElapsedTime(logger, time.Now(), "%s done: %s", r.Method, r.URL.Path)
		}

		origURI := externalURI
		forwardedURI := r.Header.Get("X-Forwarded-For")

		rec := httptest.NewRecorder()

		// first run the next handler and get results
		h.ServeHTTP(rec, r)

		// read response body and replace links
		bytes := rec.Body.Bytes()
		s := string(bytes)
		if len(s) > 0 {
			if len(forwardedURI) > 0 {
				// if both are changed (X-Forwarded-For and External uri environment variabele) use the last one
				if origURI == "http://localhost:8080/" {
					s = strings.Replace(s, "localhost", forwardedURI, -1)
				}
			}

		}
		// handle headers too...
		for k, v := range rec.HeaderMap {
			val := v[0]

			// if there is a location header and a proxy running, change
			// the location url too
			if k == "Location" && len(forwardedURI) > 0 {
				if origURI == "http://localhost:8080/" {
					logger.Debugf("proxy + location header detected. forwarded uri: %s", forwardedURI)
					// idea: run net.LookupAddr(forwardeduri) to get hostname instead of ip address?
					val = strings.Replace(val, "localhost", forwardedURI, -1)
				}
			}

			// add the header to response
			w.Header().Add(k, val)
		}

		// handle status code
		w.WriteHeader(rec.Code)

		// write modified response
		w.Write([]byte(s))
	}
	return http.HandlerFunc(fn)
}
