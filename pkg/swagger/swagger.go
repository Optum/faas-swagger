package swagger

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ghodss/yaml"
	types "github.com/openfaas/faas-provider/types"
	"github.com/pkg/errors"

	"github.com/optum/faas-swagger/pkg/auth"
)

const (
	functions_path         = "/system/functions"
	base_swagger_yaml_path = "/var/openfaas/secrets/swagger.yaml"
)

var (
	defaultStructure = `{
												   "post": {
												      "requestBody": {
												         "content": {
												            "application/json": {
												               "schema": {
												                  "type": "object"
												               }
												            }
												         }
												      },
												      "responses": {
												         "200": {
												            "description": "success"
												         }
												      }
												   }
												}`
	// EmptyResponse indicates the error that Openfaas Gateway
	// returned a empty list of functions
	EmptyResponse = errors.New("Empty response from Openfaas gateway")
)

//SwaggerConstructor holds the state of swagger constructor
type SwaggerConstructor struct {
	Gateway          string
	authPlugin       auth.OFAuth
	defaultStructure map[string]interface{}
	sYAML            map[string]interface{}
}

//Constructor returns an instance of SwaggerConstructor
func Constructor(authPlugin auth.OFAuth) *SwaggerConstructor {
	var swaggerYAML map[string]interface{}
	dat, _ := ioutil.ReadFile(base_swagger_yaml_path)
	yaml.Unmarshal(dat, &swaggerYAML)
	if swaggerYAML["paths"] == nil {
		swaggerYAML["paths"] = make(map[string]interface{})
	}

	var def map[string]interface{}
	json.Unmarshal([]byte(defaultStructure), &def)

	return &SwaggerConstructor{
		"http://gateway:8080", //overridable
		authPlugin,
		def,
		swaggerYAML,
	}
}

//GetSwaggerYAML constructs the swagger yaml for all the functions
//registerd in the openfaas gateway
func (c *SwaggerConstructor) GetSwaggerYAML() ([]byte, error) {
	paths := (c.sYAML["paths"]).(map[string]interface{})
	fndata, err := c.getFunctionsList()
	if err != nil {
		return nil, err
	}
	for _, fn := range fndata {
		anns := *fn.Annotations
		swaggerAnn := anns["swagger"]
		if swaggerAnn != "" {
			y := make(map[string]interface{})
			err := json.Unmarshal([]byte(swaggerAnn), &y)
			if err != nil {
				//ignore not well formed json and continue
				log.Printf("json not formed well %s\n", err)
				paths["/"+fn.Name] = c.defaultStructure
				continue
			}
			paths["/"+fn.Name] = y
		} else if paths["/"+fn.Name] == nil {
			paths["/"+fn.Name] = c.defaultStructure
		} else {
			continue
		}
	}
	return yaml.Marshal(c.sYAML)
}

func (c *SwaggerConstructor) getFunctionsList() ([]types.FunctionStatus, error) {
	req, err := http.NewRequest("GET", c.Gateway+functions_path, nil)
	c.authPlugin.AddAuth(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Can't connect to the given openfaas gateway")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bytesOut, _ := ioutil.ReadAll(resp.Body)

	if string(bytesOut) == "" {
		return nil, EmptyResponse
	}
	var fndata []types.FunctionStatus
	json.Unmarshal(bytesOut, &fndata)
	return fndata, nil
}
