package basic

import (
	"log"
	"net/http"	
	"io/ioutil"
)

const (
	userSecretPath     = "/var/openfaas/secrets/basic-auth-user"
	passwordSecretPath = "/var/openfaas/secrets/basic-auth-password"
)

type BasicAuth struct {
	user string
	pass string
}

func Init() *BasicAuth {
	return &BasicAuth{
		string(getBasicAuthSecret(userSecretPath)),
		string(getBasicAuthSecret(passwordSecretPath)),
	}
}

func (ba *BasicAuth) AddAuth(req *http.Request) error {
	req.SetBasicAuth(ba.user, ba.pass)
	return nil
}

func getBasicAuthSecret(path string) []byte {
	secretBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to retrieve basic auth creds, %v\n", err)
	}
	return secretBytes
}
