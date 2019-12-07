package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/ghodss/yaml"
	types "github.com/openfaas/faas-provider/types"
	"github.com/rakyll/statik/fs"
	_ "github.optum.com/Optum-Serverless/faas-swagger/statik"
)

func main() {

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	gateway, _ := os.LookupEnv("openfaas_gateway")
	u, _ := url.Parse(gateway+"/function/")

	http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", http.FileServer(statikFS)))
	http.HandleFunc("/swagger.yaml", generateSwaggerYmlHandler)

	//reverse proxy openfaas requests
	http.Handle("/", httputil.NewSingleHostReverseProxy(u))

	http.ListenAndServe(":8080", nil)

}

//generate swagger yml
func generateSwaggerYmlHandler(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{}

	filePath, filePresent := os.LookupEnv("swagger_file_path")

	if !filePresent {
		log.Fatal("Provide swagger_file_path env variable")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dat, _ := ioutil.ReadFile(filePath + "/swagger.yaml")

	yaml.Unmarshal(dat, &result)

	addPaths(result)

	out, _ := yaml.Marshal(&result)

	log.Println("successfully generated yml")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(out))

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
	filePath, _ := os.LookupEnv("swagger_file_path")
	dat, _ := ioutil.ReadFile(filePath + "/sample.yaml")
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
