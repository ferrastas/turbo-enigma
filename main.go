package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func assert(e error) {
	if e != nil {
		panic(e)
	}
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
} 

var (Client HTTPClient)

func init() {
	Client = &http.Client{}
}

func postJson(url string, message []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(message))
	assert(err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	assert(err)

	fmt.Printf("value(%v)", resp)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return nil
}

func main() {
	guardEnvVars()

	var server = fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT"))

	fmt.Println("Server listening on", server)

	http.HandleFunc("/", postOnSlack)
	http.ListenAndServe(server, nil)
}
