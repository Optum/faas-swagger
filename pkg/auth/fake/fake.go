package fake

import (
	"net/http"
)

type FakeAuth struct{}

//AddAuth adds auth in the http request header
func (a *FakeAuth) AddAuth(req *http.Request) error {
	return nil
}
