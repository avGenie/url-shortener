package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	"os"
	"strings"

	"github.com/avGenie/url-shortener/internal/app/config"
)

func main() {
	config.ParseConfig()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var data = "https://practicum.yandex.ru/"
	url, err := postRequest(client, data)
	if err != nil {
		log.Println("Post request error: : %w", err)
		panic(err)
	}

	getRequest(client, url)
}

func readFromConsole() string {
	// приглашение в консоли
	fmt.Println("Input URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	data, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	data = strings.TrimSuffix(data, "\n")

	return data
}

func postRequest(client *http.Client, data string) (string, error) {
	fmt.Println("Post request")
	url := fmt.Sprintf("http://%s", config.Config.NetAddr)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(data)))
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return "", err
	}

	log.Println("Post request output: ", string(bodyBytes))

	return path.Base(string(bodyBytes)), nil
}

func getRequest(client *http.Client, url string) {
	requestURL := fmt.Sprintf("%s/%s", config.Config.BaseURIPrefix, url)
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Request has not been created: %w", err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Request has not been sent: %w", err)
		return
	}
	defer response.Body.Close()

	log.Println("Output location: ", response.Header.Get("Location"))
}
