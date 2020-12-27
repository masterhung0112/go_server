package api

import (
	"github.com/masterhung0112/hk_server/web"
	"net/http"
)

type Context = web.Context

func (api *API) ApiHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		HandlerName:         web.GetHandlerName(h),
		RequireSession:      false,
	}
	return handler
}