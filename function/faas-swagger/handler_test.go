package function

import (
	"fmt"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/optum/faas-swagger/lib/test"
	"github.com/optum/faas-swagger/pkg/auth/fake"
	"github.com/optum/faas-swagger/pkg/swagger"
)

var _ = (func() interface{} {
	_testing = true
	return nil
}())

func TestSwaggerGeneratorFunction(t *testing.T) {
	c = &swagger.SwaggerConstructor{
		AuthPlugin:       &fake.FakeAuth{},
		DefaultStructure: swagger.DefaultStructure(),
		BaseYAML:         swagger.BaseStructure("./../swagger.yaml"),
	}

	for _, tc := range []struct {
		name        string
		funcFile    string
		swaggerFile string
	}{{
		"case1",
		"./../../lib/test/test-files/functions/func1.json",
		"./../../lib/test/test-files/swagger/swagger1.yaml",
	}, {
		"case2",
		"./../../lib/test/test-files/functions/func2.json",
		"./../../lib/test/test-files/swagger/swagger2.yaml",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, c, tc.funcFile, tc.swaggerFile)
		})
	}
}

func runTest(t *testing.T, c *swagger.SwaggerConstructor, funcFile, swaggerFile string) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := test.GetDataFromFile(funcFile)
		fmt.Fprintln(w, string(data))
	}))
	defer ts.Close()
	c.Gateway = ts.URL

	var expected, actual map[string]interface{}
	test.GetDataFromFileInFormat(swaggerFile, &expected)

	actualResp := Handle(nil)
	err := yaml.Unmarshal([]byte(actualResp), &actual)
	assert.NilError(t, err)

	assert.DeepEqual(t, &actual, &expected)
}
