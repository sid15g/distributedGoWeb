package main

/**
	Reference:
		* https://golang.org/pkg/net/http/
		* https://gowebexamples.com/http-server/
		* https://dev.to/shindakun/attempting-to-learn-go---consuming-a-rest-api-5c7g
 */

import (
	"fmt"
	log "github.com/alexcesaro/log/stdlog"
	"net/http"
	"time"
)

var (
	logger = log.GetFromFlags()
)

func main() {
	s := createServer("127.0.0.1", int16(8080))

	s.register("/check", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw,"Check confirmed...")
	})

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