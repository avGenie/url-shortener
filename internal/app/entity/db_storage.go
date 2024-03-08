package entity

import "context"

type DBStorage interface {
	Close() error
	PingDBServer(context.Context) (int, error)
}
