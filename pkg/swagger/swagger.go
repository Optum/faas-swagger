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
	AuthPlugin       auth.OFAuth
	DefaultStructure map[string]interface{}
	BaseYAML         map[string]interface{}
}

//DefaultStructure returns default swagger doc
//incase the function doesnt provide annotations
//the default structure is used
func DefaultStructure() map[string]interface{} {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(defaultStructure), &def); err != nil {
		log.Println("error getting default structure", err)
	}
	return def
}

//BaseStructure return the base swagger yaml
//to be populated with docs for each function
func BaseStructure(filePath string) map[string]interface{} {
	var base map[string]interface{}
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("error loading base structure file", err)
	}
	if err := yaml.Unmarshal(dat, &base); err != nil {
		log.Println("error getting base structure", err)
	}
	if base["paths"] == nil {
		base["paths"] = make(map[string]interface{})
	}
	return base
}

//Constructor returns an instance of SwaggerConstructor
func Constructor(gatewayURL string, authPlugin auth.OFAuth) *SwaggerConstructor {
	return &SwaggerConstructor{
		gatewayURL,
		authPlugin,
		DefaultStructure(),
		BaseStructure(base_swagger_yaml_path),
	}
}

//GetSwaggerYAML constructs the swagger yaml for all the functions
//registerd in the openfaas gateway
func (c *SwaggerConstructor) GetSwaggerYAML() ([]byte, error) {
	paths := (c.BaseYAML["paths"]).(map[string]interface{})
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
				paths["/"+fn.Name] = c.DefaultStructure
				continue
			}
			paths["/"+fn.Name] = y
		} else if paths["/"+fn.Name] == nil {
			paths["/"+fn.Name] = c.DefaultStructure
		} else {
			continue
		}
	}
	return yaml.Marshal(c.BaseYAML)
}

func (c *SwaggerConstructor) getFunctionsList() ([]types.FunctionStatus, error) {
	req, err := http.NewRequest("GET", c.Gateway+functions_path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to the given openfaas gateway")
	}
	c.AuthPlugin.AddAuth(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to the given openfaas gateway")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bytesOut, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error getting payload")
	}
	if string(bytesOut) == "" {
		return nil, EmptyResponse
	}
	var fndata []types.FunctionStatus
	json.Unmarshal(bytesOut, &fndata)
	return fndata, nil
}
