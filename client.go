package distributedGoWeb

import (
	"io/ioutil"
	"net/http"
)


type client struct {

}

func (c *client) newRequest(method string, apiUrl string) {

	req, err := http.NewRequest(method, apiUrl, nil)

	if check(err) == true {

		reqBody, err := ioutil.ReadAll(req.Body);

		if check(err) == true {

			logger.Info(reqBody);
			//TODO send request
			/*
				client := http.DefaultClient
			    resp, err := client.Do(req)
			*/
		}
	}

}//end of function
