package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	maxEncodedSize = 8
)

var (
	urls = make(map[string]string)

	ErrInvalidGivenURL = errors.New("given URL is mapped")

	EmptyURL        = "URL is empty"
	ShortURLNotInDB = "given short URL did not find in database"
)

// Processes POST request. Sends short URL in http://localhost:8080/id format.
//
// Encodes given URL using base64 encoding scheme and puts it to the URL's map.
//
// Returns 201 status code if processing was successfull, otherwise returns 400.
func PostHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		IncorrectRequestHandler(writer, req, fmt.Sprintf("cannot process URL: %s", err.Error()))
		return
	}

	if len(bodyBytes) == 0 {
		IncorrectRequestHandler(writer, req, EmptyURL)
		return
	}

	bodyString := string(bodyBytes)
	encodedURL := base64.StdEncoding.EncodeToString(bodyBytes)
	shortURL := encodedURL[:maxEncodedSize]

	urls[shortURL] = bodyString

	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, fmt.Sprintf("http://%s/%s", req.Host, shortURL))
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func GetHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")

	shortURL, err := processInputData(req.RequestURI)
	if err != nil {
		IncorrectRequestHandler(writer, req, err.Error())
		return
	}

	decodedURL, ok := urls[shortURL]
	if !ok {
		IncorrectRequestHandler(writer, req, ShortURLNotInDB)
		return
	}

	writer.Header().Set("Location", decodedURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

// Sends 400 status code with error message
func IncorrectRequestHandler(writer http.ResponseWriter, req *http.Request, message string) {
	writer.WriteHeader(http.StatusBadRequest)
	io.WriteString(writer, message)
}

// Removes "/" prefix from the input data if length >= 2
func processInputData(data string) (string, error) {
	if len(data) < 2 || !strings.HasPrefix(data, "/") {
		return "", ErrInvalidGivenURL
	}

	return data[1:], nil
}
