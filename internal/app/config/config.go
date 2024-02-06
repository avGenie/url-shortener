package config

import (
	"flag"
)

var (
	NetAddr string
	BaseURIPrefix string
)

func init() {
	flag.StringVar(&NetAddr, "a", "localhost:8080", "net address host:port")
	flag.StringVar(&BaseURIPrefix, "b", "http://localhost:8080", "base output short URL")
}
