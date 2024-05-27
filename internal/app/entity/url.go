package entity

import (
	"fmt"
	"net/url"
)

// URL Contains data related to url
type URL struct {
	Scheme string
	Host   string
	Path   string
}

// NewURL Creates URL based on input string
//
// Returns error due to parsing input string
func NewURL(inputURL string) (*URL, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}

	newURL := &URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}
	return newURL, nil
}

// ParseURL Parses URL
//
// If URL couldn't be parsed, returns nil
func ParseURL(inputURL string) (*URL, error) {
	u, err := NewURL(inputURL)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// IsValidURL Validates URL obtained from input string
func IsValidURL(inputURL string) bool {
	if len(inputURL) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(inputURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// String Implements stringer interface
func (u URL) String() string {
	s, err := url.JoinPath(u.Host, u.Path)
	if err != nil {
		return ""
	}

	if u.Scheme != "" {
		s = fmt.Sprintf("%s://%s", u.Scheme, s)
	}

	return s
}
