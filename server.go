package distributedGoWeb

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

type server struct {
	ip string
	port int16
}


func (s *server) servAt(restUri string) {

	apiUrl := "http://"+ s.ip +":"+ string(s.port) +"/"+ restUri
	go s.mapRequest(http.MethodGet, apiUrl);

}


func (s *server) mapRequest(method string, apiUrl string) {

	req, err := http.NewRequest(method, apiUrl, nil)

	if err != nil {
		panic(err)
	} else {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		} else {
			//TODO map request
			fmt.Print(reqBody)
		}
	}

}