// +build e2e

package e2e

import (
	"gotest.tools/assert"
	"testing"

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

}
