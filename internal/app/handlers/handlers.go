package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var urls = make(map[string]string)

const maxEncodedSize = 8

var ErrInvalidGivenURL = errors.New("given URL is mapped")

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandler(writer http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		IncorrectRequestHandler(writer, req, fmt.Sprintf("cannot process URL: %s", err.Error()))
		return
	}

	if len(bodyBytes) > 0 {
		bodyString := string(bodyBytes)
		encodedURL := base64.StdEncoding.EncodeToString(bodyBytes)
		shortURL := encodedURL[:maxEncodedSize]

		urls[shortURL] = bodyString

		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusCreated)
		io.WriteString(writer, fmt.Sprintf("http://%s/%s", req.Host, shortURL))
	} else {
		IncorrectRequestHandler(writer, req, "URL is empty")
	}
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func GetHandler(writer http.ResponseWriter, req *http.Request) {
	shortURL, err := processInputData(req.RequestURI)
	if err != nil {
		IncorrectRequestHandler(writer, req, err.Error())
		return
	}

	if decodedURL, ok := urls[shortURL]; ok {
		writer.Header().Set("Content-Type", "text/plain")
		writer.Header().Set("Location", decodedURL)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		IncorrectRequestHandler(writer, req, "given short URL did not find in database.")
	}
}

// Sends 400 status code with error message
func IncorrectRequestHandler(writer http.ResponseWriter, req *http.Request, message string) {
	writer.WriteHeader(http.StatusBadRequest)
	io.WriteString(writer, message)
}

// Remove "/" prefix from the input data if length >= 2
func processInputData(data string) (string, error) {
	if len(data) < 2 || !strings.HasPrefix(data, "/") {
		return "", ErrInvalidGivenURL
	}

	return data[1:], nil
}
