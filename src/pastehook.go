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
	"github.com/atotto/clipboard"
)

const (
	nilInput   = "nil"
	apiurl     = "https://pastebin.com/api/api_post.php"
	api_option = "paste"
)

func isLinesInInput() bool {
	if len(os.Args) <= 2 {
		return false
	} else {
		return true
	}
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
	return number1, number2

}

func parseFile(fileToOpen string, linesToSend string) string {
	file, err := os.Open(fileToOpen)
	red := io.Reader(file)
	var lineStart, lineEnd int64
	if linesToSend == nilInput {
		lineStart, lineEnd = 0, countLines(red)
	} else {	
		lineStart, lineEnd = splitInput(linesToSend)
	}
	if err != nil {
		panic(err)
	}
	file.Seek(0, io.SeekStart)
	neccessaryLines := scanLines(red, lineStart, lineEnd)
	file.Close()
	return neccessaryLines
}

func scanLines(red io.Reader, lineStart int64, lineEnd int64)string{
	scanner := bufio.NewScanner(red)
	scanner.Split(bufio.ScanLines)
	var buffer strings.Builder
	var line int64
	for scanner.Scan(){
		line++
		if line >= lineStart{
			buffer.WriteString(scanner.Text())
			buffer.WriteString("\n")
		}
		if line == lineEnd{
			break
		}
			
	}
	return buffer.String()
}

func countLines(red io.Reader) int64 {
	scanner := bufio.NewScanner(red)
	var lineCount int64	
	for scanner.Scan(){
		lineCount++
	}
	if err := scanner.Err(); err != nil{
		return lineCount
	}
	return lineCount 
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

func formData(text string)url.Values{
	v := url.Values{}
	v.Set("api_dev_key", getDevKey())
	v.Set("api_option", api_option)
	v.Set("api_paste_code", text)
	return v
}

func sendTextToPastehook(data url.Values)[]byte{
	request, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer([]byte(data.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return responseBody
}

func putLinkToClipboard(link string){
	err := clipboard.WriteAll(link)
	if err != nil {
		fmt.Println("Error with copying link to buffer")
		panic(err)
	}
}

func main() {
	fileToOpen, linesToSend := inputArgs()
	text := parseFile(fileToOpen, linesToSend)
	data := formData(text)
	putLinkToClipboard(string(sendTextToPastehook(data)))
}
