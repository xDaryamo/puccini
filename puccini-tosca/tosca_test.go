package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/problems"
)

func TestParse(t *testing.T) {
	testParse(t, "grammar/artifacts.yaml", nil)
	testParse(t, "grammar/attributes.yaml", nil)
	testParse(t, "grammar/data-types.yaml", nil)
	testParse(t, "grammar/descriptions.yaml", nil)
	testParse(t, "grammar/dsl-definitions.yaml", nil)
	testParse(t, "grammar/functions.yaml", nil)
	testParse(t, "grammar/inputs-and-outputs.yaml", map[string]interface{}{"ram": "1gib"})
	testParse(t, "grammar/interfaces.yaml", nil)
	testParse(t, "grammar/metadata.yaml", nil)
	testParse(t, "grammar/namespaces.yaml", nil)
	testParse(t, "grammar/policies-and-groups.yaml", nil)
	testParse(t, "grammar/requirements-and-capabilities.yaml", nil)
	testParse(t, "grammar/simple-for-nfv.yaml", nil)
	testParse(t, "grammar/source-and-target.yaml", nil)
	testParse(t, "grammar/substitution-mapping-client.yaml", nil)
	testParse(t, "grammar/substitution-mapping.yaml", nil)
	testParse(t, "grammar/unicode.yaml", nil)
	testParse(t, "grammar/workflows.yaml", nil)
	testParse(t, "javascript/constraints.yaml", nil)
	testParse(t, "javascript/functions.yaml", nil)
	testParse(t, "javascript/exec.yaml", nil)
	testParse(t, "javascript/xml.yaml", nil)
	testParse(t, "kubernetes/bookinfo/bookinfo-simple.yaml", nil)
}

var ROOT string

func init() {
	ROOT = os.Getenv("ROOT")
}

func testParse(t *testing.T, url string, inputs map[string]interface{}) {
	t.Run(url, func(t *testing.T) {
		// Running the tests in parallel is not for speed -
		// it actually allowed us to find several concurrency bugs
		t.Parallel()

		var s *normal.ServiceTemplate
		var c *clout.Clout
		var p *problems.Problems
		var err error

		if s, p, err = parser.Parse(fmt.Sprintf("%s/examples/%s", ROOT, url), inputs); err != nil {
			t.Errorf("%s\n%s", err.Error(), p)
			return
		}

		if c, err = compiler.Compile(s); err != nil {
			t.Errorf("%s\n%s", err.Error(), p)
			return
		}

		compiler.Coerce(c, p)
		if !p.Empty() {
			t.Errorf("%s", p)
		}
	})
}
