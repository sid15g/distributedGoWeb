package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type server struct {
	ip string
	port int16
	sock *http.Server
}

type HttpHandler func(http.ResponseWriter, *http.Request)
type MsgListener net.Listener

/**
 *	For the dynamic aspect, the http.Request contains all information
 *  about the request and it’s parameters. You can read GET parameters
 *  with r.URL.Query().Get("token") or POST parameters (fields from an HTML form)
 *  with r.FormValue("email").
*/
func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Welcome to my GoLang Server! ");
	logger.Info(" GET / 200 OK ")
}

func myHttpHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof(" GET %s 200 OK ", r.URL.Path )
		h(rw, r)
	}
}

func (s *server) register(path string, handler http.HandlerFunc ) {
	logger.Infof(" Registering [%s] ", path)
	http.HandleFunc(path, myHttpHandler(handler));
}


/**
 * Can also use http.ListenAndServe(":8080", nil)
 * But thats a naive implementation
 * Server.ListenAndServe() implements KeepAlive internally and
 * is a better implementation
*/
func (s *server) startWith(readTimeout int64, writeTimeout int64) error	{
	s.sock = &http.Server{
		Addr:           ":"+ strconv.Itoa(int(s.port)),
		Handler:        s,
		ReadTimeout:    time.Duration(readTimeout),
		WriteTimeout:   time.Duration(writeTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	logger.Infof(" Starting server at %s@%d ", s.ip, s.port )
	return s.sock.ListenAndServe();
}

func (s *server) start() error	{
	return s.startWith(int64(5 * time.Second), int64(15 * time.Second));
}

func (s *server) stop() error	{
	logger.Info(" Web Server Stopped..!!! ")
	return s.sock.Close()
}

func createServer(ip string, port int16) *server {

	s := &server{ip:ip, port:port}

	//Register the handler
	http.HandleFunc("/", s.ServeHTTP);
	return s;

}