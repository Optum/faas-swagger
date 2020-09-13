package auth

import (
	"net/http"
	"os"

	"github.com/optum/faas-swagger/pkg/auth/basic"
)

type OFAuth interface {
	AddAuth(req *http.Request) error
}

func GetAuthPlugin() OFAuth {
	switch os.Getenv("AUTH_TYPE") {
	case "OIDC":
		//TBI
	default:
		return basic.Init()
	}
	return nil
}
