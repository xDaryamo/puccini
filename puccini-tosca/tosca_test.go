package main

import (
	"fmt"
	"os"
	"testing"

	problemspkg "github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
)

func TestParse(t *testing.T) {
	testCompile(t, "tosca/artifacts.yaml", nil)
	testCompile(t, "tosca/attributes.yaml", nil)
	testCompile(t, "tosca/copy.yaml", nil)
	testCompile(t, "tosca/data-types.yaml", nil)
	testCompile(t, "tosca/descriptions.yaml", nil)
	testCompile(t, "tosca/dsl-definitions.yaml", nil)
	testCompile(t, "tosca/functions.yaml", nil)
	testCompile(t, "tosca/inputs-and-outputs.yaml", map[string]interface{}{"ram": "1gib"})
	testCompile(t, "tosca/interfaces.yaml", nil)
	testCompile(t, "tosca/metadata.yaml", nil)
	testCompile(t, "tosca/namespaces.yaml", nil)
	testCompile(t, "tosca/policies-and-groups.yaml", nil)
	testCompile(t, "tosca/requirements-and-capabilities.yaml", nil)
	testCompile(t, "tosca/simple-for-nfv.yaml", nil)
	testCompile(t, "tosca/source-and-target.yaml", nil)
	testCompile(t, "tosca/substitution-mapping-client.yaml", nil)
	testCompile(t, "tosca/substitution-mapping.yaml", nil)
	testCompile(t, "tosca/unicode.yaml", nil)
	testCompile(t, "tosca/workflows.yaml", nil)
	testCompile(t, "tosca/legacy/tosca_1_0.yaml", nil)
	testCompile(t, "tosca/legacy/tosca_1_1.yaml", nil)
	testCompile(t, "tosca/legacy/tosca_1_2.yaml", nil)
	testCompile(t, "tosca/future/tosca_2_0.yaml", nil)
	testCompile(t, "javascript/artifacts.yaml", nil)
	testCompile(t, "javascript/constraints.yaml", nil)
	testCompile(t, "javascript/define.yaml", nil)
	testCompile(t, "javascript/exec.yaml", nil)
	testCompile(t, "javascript/functions.yaml", nil)
	testCompile(t, "kubernetes/bookinfo/bookinfo-simple.yaml", nil)
	testCompile(t, "openstack/hello-world.yaml", nil)
	testCompile(t, "bpmn/open-loop.yaml", nil)
	testCompile(t, "cloudify/advanced-blueprint-example.yaml", map[string]interface{}{
		"host_ip":                "1.2.3.4",
		"agent_user":             "my_user",
		"agent_private_key_path": "my_key",
	})
	testCompile(t, "cloudify/example.yaml", nil)
	testCompile(t, "hot/hello-world.yaml", map[string]interface{}{
		"username": "test",
	})
}

var ROOT string

func init() {
	ROOT = os.Getenv("ROOT")
}

func testCompile(t *testing.T, url string, inputs map[string]interface{}) {
	t.Run(url, func(t *testing.T) {
		// Running the tests in parallel is not just for speed;
		// it actually helps us to find concurrency bugs
		t.Parallel()

		var serviceTemplate *normal.ServiceTemplate
		var clout *cloutpkg.Clout
		var problems *problemspkg.Problems
		var err error

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		url_, err := urlpkg.NewURL(fmt.Sprintf("%s/examples/%s", ROOT, url), urlContext)
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}

		if _, serviceTemplate, problems, err = parser.Parse(url_, nil, inputs); err != nil {
			t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
			return
		}

		if clout, err = compiler.Compile(serviceTemplate, true); err != nil {
			t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
			return
		}

		compiler.Resolve(clout, problems, urlContext, true, "yaml", false, false, true)
		if !problems.Empty() {
			t.Errorf("%s", problems.ToString(true))
			return
		}

		compiler.Coerce(clout, problems, urlContext, true, "yaml", false, false, true)
		if !problems.Empty() {
			t.Errorf("%s", problems.ToString(true))
			return
		}
	})
}
