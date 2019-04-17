package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/**
 *	Reference:
 *		https://gobyexample.com/worker-pools
 *		https://stackoverflow.com/questions/18293133/using-goroutines-for-background-work-inside-an-http-handler
 */

type Method string

const (
	GET Method = "GET"
    POST Method = "POST"
)

type HttpHandlerAdapter struct {
	path string
	method Method
	defaultStatusCode int
	handler http.HandlerFunc
}

type server struct {
	ip string
	port int16
	sock *http.Server
	rmap map[string]HttpHandlerAdapter
}

func (s *server) status404NotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound);
	fmt.Fprintf(rw, " Error 404 Page Not Found ")
	logger.Infof(" %s %s 404 NotFound ", r.Method, r.URL.Path )
}

/**
 *	For the dynamic aspect, the http.Request contains all information
 *  about the request and it’s parameters. You can read GET parameters
 *  with r.URL.Query().Get("token") or POST parameters (fields from an HTML form)
 *  with r.FormValue("email").
*/
func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	mthd := r.Method

	if r.URL.Path == "/" {
		rw.WriteHeader(http.StatusOK);
		fmt.Fprint(rw, " Welcome to GoLang Web Server! ");
		logger.Infof(" %s %s 200 OK ", mthd, r.URL.Path )
	}else if val, ok := s.rmap[r.URL.Path]; ok {
		if string(val.method) == mthd {
			rw.WriteHeader(val.defaultStatusCode);
			val.handler(rw, r)
		} else {
			s.status404NotFound(rw, r)
		}
	} else {
		s.status404NotFound(rw, r)
	}

}

/**
 *  Go multiplexes goroutines onto the available threads which is determined by the
 *  GOMAXPROCS environment setting. As a result if this is set to 1 then a single
 *  goroutine can hog the single thread. Go has available to it until it yields control back to the Go runtime.
 *  More than likely doSomeBackgroundWork is hogging all the time on a single thread which is
 *  preventing the http handler from getting scheduled.
 *
 *  First, as a general rule when using goroutines, you should set GOMAXPROCS to the number of
 *  CPUs your system has or to whichever is bigger.
 *
 */
func createHttpHandler(h http.HandlerFunc, method Method, defaultStatusCode int) HttpHandlerAdapter {

	adapter := HttpHandlerAdapter {
		handler: func (rw http.ResponseWriter, r *http.Request) {
			logger.Infof(" %s %s 200 OK ", r.Method, r.URL.Path )
			h(rw, r)
		},
		defaultStatusCode: defaultStatusCode,
		method: method};

	return adapter;
}

func defaultHttpHandler(h http.HandlerFunc) HttpHandlerAdapter {
	return createHttpHandler(h, GET, http.StatusOK)
}


func (s *server) register(path string, handler http.HandlerFunc) *server {
	logger.Infof(" Registering [%s] ", path)
	h := defaultHttpHandler(handler);
	h.path = path
	s.rmap[path] = h
	return s
}


/**
 * 	Reference:
 * 		https://www.alexedwards.net/blog/serving-static-sites-with-go
*/
func (s *server) serveStatic(path string) *server {
	//TODO code
	return s;
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

	logger.Infof(" Go WebServer started at Port %d ", s.port )
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
	s.rmap = make(map[string]HttpHandlerAdapter)
	return s;
}