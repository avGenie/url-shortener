package main

import (
	"bufio"
	"fmt"
	"path"

	"os"
	"strings"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

func main() {
	cli := gentleman.New()
	cli.URL("http://localhost:8080/")

	// data := readFromConsole()
	var data = "https://practicum.yandex.ru/"
	url, err := postRequest(cli, data)
	if err != nil {
		panic(err)
	}

	getRequest(cli, url)
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

func postRequest(cli *gentleman.Client, data string) (string, error) {
	req := cli.Request()
	req.Method("POST")
	req.Use(body.String(data))

	// Perform the request
	res, err := req.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return "", fmt.Errorf("request error")
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return "", fmt.Errorf("request error")
	}

	output := res.String()

	return path.Base(output), nil
}

func getRequest(cli *gentleman.Client, url string) {
	fmt.Printf("Input url: %s\n", url)

	req := cli.Request()
	req.Method("GET")
	req.Path(fmt.Sprintf("/%s", url))

	// Perform the request
	res, err := req.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return
	}

	// Reads the whole body and returns it as string
	fmt.Printf("Body: %s", res.Header.Get("Location"))
}
