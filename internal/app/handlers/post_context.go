package handlers

import "net/http"

type PostContext struct {
	baseURIPrefix string
	handle        func(string, http.ResponseWriter, *http.Request)
}

func (ctx *PostContext) Handle() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx.handle(ctx.baseURIPrefix, writer, req)
	}
}
