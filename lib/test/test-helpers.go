package test

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
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
