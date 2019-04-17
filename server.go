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
 */
type HttpHandlerAdapter struct {
	path string
	handler http.HandlerFunc
}

type server struct {
	ip string
	port int16
	sock *http.Server
	rmap map[string]HttpHandlerAdapter
}

/**
 *	For the dynamic aspect, the http.Request contains all information
 *  about the request and it’s parameters. You can read GET parameters
 *  with r.URL.Query().Get("token") or POST parameters (fields from an HTML form)
 *  with r.FormValue("email").
*/
func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		logger.Infof(" GET %s 200 OK ", r.URL.Path )
		rw.WriteHeader(http.StatusOK);
		fmt.Fprint(rw, "Welcome to my GoLang Server! ");
	}else if val, ok := s.rmap[r.URL.Path]; ok {
		rw.WriteHeader(http.StatusOK);
		val.handler(rw, r)
	} else {
		rw.WriteHeader(http.StatusNotFound);
		fmt.Fprintf(rw, " Error 404 Page Not Found ")
		logger.Infof(" GET %s 404 NotFound ", r.URL.Path )
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
func myHttpHandler(h http.HandlerFunc) HttpHandlerAdapter {

	adapter := HttpHandlerAdapter {
		handler: func (rw http.ResponseWriter, r *http.Request){
			logger.Infof(" GET %s 200 OK ", r.URL.Path )
			h(rw, r)
		}};

	return adapter;
}

func (s *server) register(path string, handler http.HandlerFunc ) {
	logger.Infof(" Registering [%s] ", path)
	h := myHttpHandler(handler);
	h.path = path
	s.rmap[path] = h
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

	//Register the handler
	s.register("/", myHttpHandler(s.ServeHTTP).handler);
	return s;

}