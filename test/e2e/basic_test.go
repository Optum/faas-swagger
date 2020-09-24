// +build e2e

package e2e

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/optum/faas-swagger/lib/test"
)

func TestBasicWorkflow(t *testing.T) {
	t.Parallel()
	it := NewTest()

	t.Log("call go hw through swagger proxy")
	it.callHW(t)

	t.Log("verify if go hw function is in swagger yaml")
	it.checkSwaggerYamlForHW(t)
}

func (it *FSTest) callHW(t *testing.T) {
	bytesOut, err := test.InvokeHTTP(it.SwaggerAddr + GO_HW_FUNCTION)
	assert.NilError(t, err)
	assert.Equal(t, GO_HW_REPONSE, string(bytesOut))
}

func (it *FSTest) checkSwaggerYamlForHW(t *testing.T) {
	bytesOut, err := test.InvokeHTTP(it.SwaggerAddr + GO_SWAGGER_FUNCTION)
	fmt.Println(string(bytesOut))
	assert.NilError(t, err)
	var actual, want map[string]interface{}
	yaml.Unmarshal(bytesOut, &actual)
	paths := (actual["paths"]).(map[string]interface{})
	gohwDoc := paths["/go-hw"]
	json.Unmarshal([]byte(GO_HW_SWAGGER_DOC),&want)
	wantDoc := paths["/go-hw"]
	assert.DeepEqual(t, &wantDoc, &gohwDoc)
}
