package entity

import "context"

type Storage interface {
	Close() Response
	PingServer(ctx context.Context) Response
	AddURL(ctx context.Context, key, value URL) Response
	GetURL(ctx context.Context, key URL) URLResponse
}
