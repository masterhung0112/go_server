package api

import (
	"github.com/masterhung0112/go_server/model"
  "github.com/gorilla/mux"
)

type ApiRoutes struct {
  Root    *mux.Router // ''
  ApiRoot *mux.Router // 'api/v1'

  Users          *mux.Router // 'api/v1/users'
	User           *mux.Router // 'api/v1/users/{user_id:[A-Za-z0-9]+}'
}

type API struct {
  BaseRoutes    *ApiRoutes
}

func ApiInit(root *mux.Router) {
  api := &API {
    BaseRoutes: &ApiRoutes{},
  }

  api.BaseRoutes.Root = root
  api.BaseRoutes.ApiRoot = root.PathPrefix(model.API_URL_SUFFIX).Subrouter()

  api.BaseRoutes.Users = api.BaseRoutes.ApiRoot.PathPrefix("/users").Subrouter()
  api.BaseRoutes.User = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-za-z0-9]+}").Subrouter()
  api.InitUser()
}