package main

/**
	Reference:
		* https://golang.org/pkg/net/http/
		* https://gowebexamples.com/http-server/
		* https://dev.to/shindakun/attempting-to-learn-go---consuming-a-rest-api-5c7g
 */

import (
	log "github.com/alexcesaro/log/stdlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger = log.GetFromFlags()
)

func main() {

	s := createServer("127.0.0.1", int16(8080))
	go shutdown(s);

	t := template{}
	s.register("/check", t.printOnly("Check confirmed..."),
	).register("/check/sid", t.printOnly("Sid confirmed..."),
	);

	if check( s.start() ) == false {
		logger.Error("Unable to start Go Web Server ")
	}

}

/* Checks for error and prints if found */
func check(err error) bool {
	if err != nil {
		logger.Alert(err)
		return false
	}
	return true
}

func shutdown(s *server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	logger.Info(" Received", sig, "| Gracefully Shutting down the web-server in 5 seconds ")
	time.Sleep(5 * time.Second)
	s.stop()
	logger.Close()
	os.Exit(0)
}