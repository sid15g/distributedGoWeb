package main

/**
	Reference:
		* https://golang.org/pkg/net/http/
		* https://gowebexamples.com/http-server/
		* https://dev.to/shindakun/attempting-to-learn-go---consuming-a-rest-api-5c7g
 */

import (
	log "github.com/alexcesaro/log/stdlog"
	"time"
)

var (
	logger = log.GetFromFlags()
)

func main() {

	s := createServer("127.0.0.1", int16(8080))
	t := template{}

	s.register("/check", t.printOnly("Check confirmed..."),
	).register("/check/sid", t.printOnly("Sid confirmed..."),
	);

	if check( s.start() ) {
		defer s.stop()
		defer logger.Close()
		time.Sleep(time.Duration(20 * time.Second))
	} else {
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