package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	apiurl     = "https://pastebin.com/api/api_post.php"
	api_option = "paste"
)

func readFile() string {
	openFile := os.Args[1]
	file, err := os.Open(openFile)
	if err != nil {
		panic(err)
	}
	byt, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(byt)
}

func getDevKey() string {
	le, found := os.LookupEnv("API_DEV_KEY")
	if (found != true) || le == "" {
		fmt.Println("API_DEV_KEY is not set and not found in envs")
		panic(found)
	}
	api_dev_key := os.Getenv("API_DEV_KEY")
	return api_dev_key
}

func main() {
	v := url.Values{}
	v.Set("api_dev_key", getDevKey())
	v.Set("api_option", api_option)
	v.Set("api_paste_code", readFile())
	//jsonBody := []byte(`{"api_dev_key":"IEOfIviozs3g5uUW5kyhHEQo8Gmfe5p2"}`)
	request, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer([]byte(v.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println(request.Body)
	if err != nil {
		panic(err)
	}
	//response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	//responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", responseBody)

}
