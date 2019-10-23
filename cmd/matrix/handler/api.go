package handler

import "github.com/go-app-cloud/goapp"

type API interface {
	Handler(party goapp.Party, dbEngine goapp.Engine)
}
