package handlers

import (
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

type GetDBContext struct {
	db     entity.DBStorage
	handle func(entity.DBStorage, http.ResponseWriter, *http.Request)
}

func (ctx *GetDBContext) Handle() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx.handle(ctx.db, writer, req)
	}
}

func NewGetDBPingContext(db entity.DBStorage) GetDBContext {
	return GetDBContext{
		db:     db,
		handle: GetPingDB,
	}
}
