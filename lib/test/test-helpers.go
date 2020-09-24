package test

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDataFromFile(pathKey string) []byte {
	byteValue, err := ioutil.ReadFile(pathKey)
	if err != nil {
		log.Println("error getting file", err)
	}
	return byteValue
}

func GetDataFromFileInFormat(pathKey string, spec interface{}) {
	byteValue := GetDataFromFile(pathKey)
	err := yaml.Unmarshal(byteValue, spec)
	if err != nil {
		log.Println("error marshalling data into given spec", err)
	}
}

func InvokeHTTP(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return ioutil.ReadAll(resp.Body)
}
