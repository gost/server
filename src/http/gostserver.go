package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"fmt"

	"github.com/geodan/gost/src/sensorthings/models"
)

// Server interface for starting and stopping the HTTP server
type Server interface {
	Start()
	Stop()
}

// GostServer is the type that contains all of the relevant information to set
// up the GOST HTTP Server
type GostServer struct {
	host string      // Hostname for example "localhost" or "192.168.1.14"
	port int         // Portnumber where you want to run your http server on
	api  *models.API // Sensorthings api to interact with from the HttpServer
}

// CreateServer initialises a new GOST HTTPServer based on the given parameters
func CreateServer(host string, port int, api *models.API) Server {
	return &GostServer{
		host: host,
		port: port,
		api:  api,
	}
}

// Start command to start the GOST HTTPServer
func (s *GostServer) Start() {
	log.Printf("Started GOST HTTP Server on %v:%v", s.host, s.port)
	router := CreateRouter(s.api)
	httpError := http.ListenAndServe(s.host+":"+strconv.Itoa(s.port), s.LowerCaseURI(router))

	if httpError != nil {
		log.Fatal(httpError)
		return
	}
}

// Stop command to stop the GOST HTTP server, currently not supported
func (s *GostServer) Stop() {

}

// LowerCaseURI is a middleware function that lower cases the url path
func (s *GostServer) LowerCaseURI(h http.Handler) http.Handler {
	api := *s.api

	fn := func(w http.ResponseWriter, r *http.Request) {
		lowerCasePath := strings.ToLower(r.URL.Path)
		split := strings.Split(lowerCasePath, "/")

		for _, split := range split {
			found := false
			for _, a := range api.acceptedPaths {
				if strings.HasPrefix(split, a) {
					found = true
				}
			}

			if !found {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusNotFound)
				w.Write(byte[0])
				return
			}
		}

		r.URL.Path = lowerCasePath
		fmt.Println(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
