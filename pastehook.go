package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	apiurl     = "https://pastebin.com/api/api_post.php"
	api_option = "paste"
)


func inputArgs() (string, string) {
	fileToOpen := os.Args[1]
	LinesToSend := os.Args[2]
	return fileToOpen, LinesToSend
}

func splitInput(linesToSend string) (int64, int64) {
	strSlice := strings.Split(linesToSend, "-")
	number1, err := strconv.ParseInt(strSlice[0], 10, 64)
	if err != nil {
		panic(err)
	}
	number2, err := strconv.ParseInt(strSlice[1], 10, 64)
	if err != nil {
		panic(err)
	}
	if (number2 - number1) < 0{
		panic("second value should be bigger than first")
	}
	fmt.Println(number1, number2)
	return number1, number2

}

func readFile(fileToOpen string, linesToSend lineNumbersToSend) string {
	file, err := os.Open(fileToOpen)
	if err != nil {
		panic(err)
	}
	red := io.Reader(file)
	scanner := bufio.NewScanner(red)
	scanner.Split(bufio.ScanLines)
	var lns []string
	for scanner.Scan(){
		lns = append(lns, scanner.Text())
	}
	file.Close()
	for position, line := range lns {
		fmt.Println(position, line)
	}
	return "ok"
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
	fileToOpen, linesToSend := inputArgs()
	number1, number2 := splitInput(linesToSend)
	v := url.Values{}
	v.Set("api_dev_key", getDevKey())
	v.Set("api_option", api_option)
	v.Set("api_paste_code", readFile(fileToOpen, number1))
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
