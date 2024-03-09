package entity

import "context"

type Storage interface {
	Close() error
	PingDBServer(context.Context) (int, error)
}
