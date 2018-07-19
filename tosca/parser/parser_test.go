package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/problems"
)

func TestParse(t *testing.T) {
	testParse(t, "tosca-grammar/artifacts.yaml", nil)
	testParse(t, "tosca-grammar/attributes.yaml", nil)
	testParse(t, "tosca-grammar/data-types.yaml", nil)
	testParse(t, "tosca-grammar/descriptions.yaml", nil)
	testParse(t, "tosca-grammar/dsl-definitions.yaml", nil)
	testParse(t, "tosca-grammar/functions-custom.yaml", nil)
	testParse(t, "tosca-grammar/functions.yaml", nil)
	testParse(t, "tosca-grammar/inputs-and-outputs.yaml", map[string]interface{}{"ram": "1gib"})
	testParse(t, "tosca-grammar/interfaces.yaml", nil)
	testParse(t, "tosca-grammar/metadata.yaml", nil)
	testParse(t, "tosca-grammar/namespaces.yaml", nil)
	testParse(t, "tosca-grammar/policies-and-groups.yaml", nil)
	testParse(t, "tosca-grammar/requirements-and-capabilities.yaml", nil)
	testParse(t, "tosca-grammar/simple-for-nfv.yaml", nil)
	testParse(t, "tosca-grammar/source-and-target.yaml", nil)
	testParse(t, "tosca-grammar/substitution-mapping-client.yaml", nil)
	testParse(t, "tosca-grammar/substitution-mapping.yaml", nil)
	testParse(t, "tosca-grammar/unicode.yaml", nil)
	testParse(t, "tosca-grammar/workflows.yaml", nil)
	testParse(t, "js/generate.yaml", nil)
	testParse(t, "js/generate-xml.yaml", nil)
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

		if s, p, err = Parse(fmt.Sprintf("%s/examples/%s", ROOT, url), inputs); err != nil {
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
