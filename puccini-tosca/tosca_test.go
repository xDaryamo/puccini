package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/problems"
)

func TestParse(t *testing.T) {
	testCompile(t, "tosca/artifacts.yaml", nil)
	testCompile(t, "tosca/attributes.yaml", nil)
	testCompile(t, "tosca/data-types.yaml", nil)
	testCompile(t, "tosca/descriptions.yaml", nil)
	testCompile(t, "tosca/dsl-definitions.yaml", nil)
	testCompile(t, "tosca/functions.yaml", nil)
	testCompile(t, "tosca/inputs-and-outputs.yaml", ard.Map{"ram": "1gib"})
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
	testCompile(t, "javascript/constraints.yaml", nil)
	testCompile(t, "javascript/exec.yaml", nil)
	testCompile(t, "javascript/functions.yaml", nil)
	testCompile(t, "kubernetes/bookinfo/bookinfo-simple.yaml", nil)
	testCompile(t, "openstack/hello-world.yaml", nil)
	testCompile(t, "bpmn/open-loop.yaml", nil)
	testCompile(t, "cloudify/advanced-blueprint-example.yaml", ard.Map{
		"host_ip":                "1.2.3.4",
		"agent_user":             "my_user",
		"agent_private_key_path": "my_key",
	})
	testCompile(t, "cloudify/example.yaml", nil)
	testCompile(t, "hot/hello-world.yaml", ard.Map{
		"username": "test",
	})
}

var ROOT string

func init() {
	ROOT = os.Getenv("ROOT")
}

func testCompile(t *testing.T, url string, inputs ard.Map) {
	t.Run(url, func(t *testing.T) {
		// Running the tests in parallel is not for speed;
		// it actually allowed us to find several concurrency bugs
		t.Parallel()

		var s *normal.ServiceTemplate
		var c *clout.Clout
		var p *problems.Problems
		var err error

		if s, p, err = parser.Parse(fmt.Sprintf("%s/examples/%s", ROOT, url), nil, inputs); err != nil {
			t.Errorf("%s\n%s", err.Error(), p)
			return
		}

		if c, err = compiler.Compile(s); err != nil {
			t.Errorf("%s\n%s", err.Error(), p)
			return
		}

		compiler.Resolve(c, p, "yaml", true)
		if !p.Empty() {
			t.Errorf("%s", p)
			return
		}

		compiler.Coerce(c, p, "yaml", true)
		if !p.Empty() {
			t.Errorf("%s", p)
			return
		}
	})
}
