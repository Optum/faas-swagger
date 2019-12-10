package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	types "github.com/openfaas/faas-provider/types"
)

// Handle a serverless request
func Handle(req []byte) string {
	var result map[string]interface{}

	dat, _ := ioutil.ReadFile("/var/openfaas/secrets/swagger.yaml")

	yaml.Unmarshal(dat, &result)

	addPaths(result)

	out, _ := yaml.Marshal(&result)

	return string(out)
}

//add paths to swagger yml
func addPaths(result map[string]interface{}) {
	pathsI := result["paths"]

	paths, _ := pathsI.(map[string]interface{})

	// initialize for first time
	if paths == nil {
		paths = make(map[string]interface{})
		result["paths"] = paths
	}

	var sample map[string]interface{}
	dat, _ := ioutil.ReadFile("/var/openfaas/secrets/sample.yaml")
	yaml.Unmarshal(dat, &sample)

	functions := getFunctionsList()

	if string(functions) == "" {
		return
	}

	var fndata []types.FunctionStatus
	json.Unmarshal(functions, &fndata)
	for _, fn := range fndata {
		ann := *fn.Annotations
		swagger := ann["swagger"]
		if swagger != "" {
			y := make(map[string]interface{})
			err := json.Unmarshal([]byte(swagger), &y)
			if err != nil {
				log.Printf("json not formed well %s", err)
				paths["/"+fn.Name] = sample
				continue
			}
			paths["/"+fn.Name] = y
		} else if paths["/"+fn.Name] == nil {
			paths["/"+fn.Name] = sample
		} else {
			continue
		}
	}
}

func getFunctionsList() []byte {
	of, filePresent := os.LookupEnv("openfaas_gateway")
	if !filePresent {
		log.Fatal("Provide openfaas_gateway env variable")
		return []byte("")
	}
	resp, err := http.Get(of + "/system/functions")
	if err != nil {
		log.Fatal("Can't connect to openfaas", err)
		return []byte("")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bytesOut, _ := ioutil.ReadAll(resp.Body)
	return bytesOut
}
