package e2e

import (
	"flag"
)

type FSTest struct {
	GatewayAddr string
	SwaggerAddr string
}

var (
	GatewayAddr = flag.String("gateway_addr", "http://127.0.0.1:8080", "The openfaas gateway addr")
	SwaggerAddr = flag.String("swagger_addr", "http://127.0.0.1:8081", "Swagger ui address")
)

const (
	GO_HW_FUNCTION      = "/go-hw"
	GO_SWAGGER_FUNCTION = "/faas-swagger-fn"
	GO_HW_REPONSE       = "Hello, Go. You said: \n"
	GO_HW_SWAGGER_DOC   = `{
												"/go-hw": {
													"get": {
														"summary": "Test function returns greetings",
														"responses": {
															"200": {
																"description": "Return greetings",
																"content": {
																	"text/plain": {
																		"schema": {
																			"type": "string",
																			"example": "hi"
																		}
																	}
																}
															},
															"401": {
																"description": "Authorization information is missing or invalid."
															},
															"404": {
																"description": "Requested resouce not found."
															},
															"405": {
																"description": "Method not allowed. Check Authorization"
															},
															"5XX": {
																"description": "Unexpected error."
															}
														}
													}
												}
											}`
)

// NewTest creates a new FSTest object
func NewTest() *FSTest {
	flag.Parse()
	return &FSTest{
		*GatewayAddr,
		*SwaggerAddr,
	}
}
