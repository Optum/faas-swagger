package function

import (
	"fmt"
	"log"

	"github.com/optum/faas-swagger/pkg/auth"
	"github.com/optum/faas-swagger/pkg/swagger"
)

var c *swagger.SwaggerConstructor

const (
	gateway_url = "http://gateway.openfaas:8080"
)

func init() {
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
