package main

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
)

func postJson(url string, json []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return err
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return nil
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println("request Body:", string(body))

	var url = os.Getenv("SLACK_WEBHOOK_URL")
	var json = []byte(`{"text":"` + html.EscapeString(string(body)) + `"}`)

	err := postJson(url, json)
	if err != nil {
		fmt.Fprintf(writer, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(writer, "OK")
}

func main() {
	var server = fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT"))

	fmt.Println("Server listening on", server)

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(server, nil)
}
