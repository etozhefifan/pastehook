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
	nilInput    = "nil"
	apiurlLogin = "https://pastebin.com/api/api_login.php"
	apiurlPaste = "https://pastebin.com/api/api_post.php"
	api_option  = "paste"
)

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
	if (number2 - number1) < 0 {
		panic("second value should be bigger than first.")
	}
	if number2-number1 == 0 {
		panic("The difference in values should not be equal to zero.")
	}
	return number1, number2
}

func scanLines(red io.Reader, lineStart int64, lineEnd int64) string {
	scanner := bufio.NewScanner(red)
	scanner.Split(bufio.ScanLines)
	var buffer strings.Builder
	var line int64
	for scanner.Scan() {
		line++
		if line >= lineStart {
			buffer.WriteString(scanner.Text())
			buffer.WriteString("\n")
		}
		if line == lineEnd-1 {
			break
		}
	}
	return buffer.String()
}

func countLines(red io.Reader) int64 {
	scanner := bufio.NewScanner(red)
	var lineCount int64
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return lineCount
	}
	return lineCount
}

func checkAndGetEnv(env string) string {
	le, found := os.LookupEnv(env)
	if (found != true) || le == "" {
		fmt.Println(env, "is not set and was not found in envs")
		panic("Env not found")
	}
	resEnv := os.Getenv(env)
	return resEnv
}


func checkForUsernameAndPassword() bool {
	leu, found := os.LookupEnv("API_USER_NAME")
	if (found != true) || leu == "" {
		return false
	}
	lep, found := os.LookupEnv("API_USER_NAME")
	if (found != true) || lep == "" {
		return false
	}
	return true
}

func getUserSession(API_DEV_KEY string, API_USER_NAME string, API_USER_PASSWORD string) string {
	v := url.Values{}
	v.Set("api_dev_key", API_DEV_KEY)
	v.Set("api_user_name", API_USER_NAME)
	v.Set("api_user_password", API_USER_PASSWORD)
	request, err := http.NewRequest(http.MethodPost, apiurlLogin, bytes.NewBuffer([]byte(v.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return string(responseBody)
}

func formData(text string) url.Values {
	v := url.Values{}
	v.Set("api_dev_key", checkAndGetEnv("API_DEV_KEY"))
	v.Set("api_option", api_option)
	v.Set("api_paste_code", text)
	if checkForUsernameAndPassword() == true {
		v.Set("api_user_key", getUserSession(os.Getenv("API_DEV_KEY"), os.Getenv("API_USER_NAME"), os.Getenv("API_USER_PASSWORD")))
	}
	return v
}

func sendTextToPastehook(data url.Values) []byte {
	request, err := http.NewRequest(http.MethodPost, apiurlPaste, bytes.NewBuffer([]byte(data.Encode())))
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

func putLinkToClipboard(link string) {
	err := clipboard.WriteAll(link)
	if err != nil {
		fmt.Println("Error with copying link to buffer")
		panic(err)
	}
	fmt.Println("Link copied to clipboard")
}

func parseFile(fileToOpen string, linesToSend string) string {
	file, err := os.Open(fileToOpen)
	red := io.Reader(file)
	var lineStart, lineEnd int64
	if linesToSend == "" {
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

func main() {
	fileToOpen, linesToSend := inputArgs()
	if fileToOpen == "" {
		panic("Specify the file using flag -f")
	}
	text := parseFile(fileToOpen, linesToSend)
	data := formData(text)
	putLinkToClipboard(string(sendTextToPastehook(data)))
  }
