package function

import (
	"fmt"
	"log"
	"os"

	"github.com/optum/faas-swagger/pkg/auth"
	"github.com/optum/faas-swagger/pkg/swagger"
)

var (
	c        *swagger.SwaggerConstructor
	_testing = false
)

const (
	default_gateway_url = "http://gateway.openfaas:8080"
)

func init() {
	if _testing {
		return
	}
	gateway_url, present := os.LookupEnv("GATEWAY_URL")
	if !present {
		gateway_url = default_gateway_url
	}
	c = swagger.Constructor(gateway_url, auth.GetAuthPlugin())
}

// Handle a serverless request
func Handle(req []byte) string {
	result, err := c.GetSwaggerYAML()
	if err != nil {
		msg := fmt.Sprintf("error constructing swagger yaml %v", err)
		log.Println(msg)
		return msg
	}
	return string(result)
}
