package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

func TestFiles(t *testing.T) {

	fileInfo, _ := ioutil.ReadDir("./test-files/files/in")

	for _, file := range fileInfo {
		if strings.HasPrefix(file.Name(), "test") {
			runTest(t, strings.Split(file.Name(), ".")[0])
		}
	}
}

func runTest(t *testing.T, fileName string) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonRespFile, _ := os.Open("./test-files/files/in/" + fileName + ".json")
		byteValue, _ := ioutil.ReadAll(jsonRespFile)
		fmt.Fprintln(w, string(byteValue))
	}))
	defer ts.Close()

	os.Setenv("swagger_file_path", "./test-files")
	os.Setenv("openfaas_gateway", ts.URL)

	expectedRespFile, _ := os.Open("./test-files/files/out-files/" + fileName + ".yaml")
	expectedResp, _ := ioutil.ReadAll(expectedRespFile)

	// Fails on invalid request
	req, err := http.NewRequest("GET", "/swagger.yaml", strings.NewReader(""))

	checkError(err, t)

	resp := httptest.NewRecorder()

	http.HandlerFunc(generateSwaggerYmlHandler).
		ServeHTTP(resp, req)

	respBytes, _ := ioutil.ReadAll(resp.Body)

	var result, expected map[string]interface{}
	yaml.Unmarshal(respBytes, &result)
	yaml.Unmarshal(expectedResp, &expected)

	if diff := cmp.Diff(result, expected); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
