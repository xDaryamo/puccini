package main

import (
	"fmt"
	"os"
	"testing"

	problemspkg "github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"

	_ "github.com/tliron/kutil/logging/simple"
)

func TestParse(t *testing.T) {
	context := NewContext(t)
	defer context.urlContext.Release()

	context.compile("tosca/artifacts.yaml", nil)
	context.compile("tosca/attributes.yaml", nil)
	context.compile("tosca/copy.yaml", nil)
	context.compile("tosca/data-types.yaml", nil)
	context.compile("tosca/descriptions.yaml", nil)
	context.compile("tosca/dsl-definitions.yaml", nil)
	context.compile("tosca/functions.yaml", nil)
	context.compile("tosca/inputs-and-outputs.yaml", map[string]interface{}{"ram": "1gib"})
	context.compile("tosca/interfaces.yaml", nil)
	context.compile("tosca/metadata.yaml", nil)
	context.compile("tosca/namespaces.yaml", nil)
	context.compile("tosca/policies-and-groups.yaml", nil)
	context.compile("tosca/requirements-and-capabilities.yaml", nil)
	context.compile("tosca/simple-for-nfv.yaml", nil)
	context.compile("tosca/source-and-target.yaml", nil)
	context.compile("tosca/substitution-mapping-client.yaml", nil)
	context.compile("tosca/substitution-mapping.yaml", nil)
	context.compile("tosca/unicode.yaml", nil)
	context.compile("tosca/workflows.yaml", nil)
	context.compile("tosca/legacy/tosca_1_0.yaml", nil)
	context.compile("tosca/legacy/tosca_1_1.yaml", nil)
	context.compile("tosca/legacy/tosca_1_2.yaml", nil)
	context.compile("tosca/future/tosca_2_0.yaml", nil)
	context.compile("javascript/artifacts.yaml", nil)
	context.compile("javascript/constraints.yaml", nil)
	context.compile("javascript/converters.yaml", nil)
	context.compile("javascript/define.yaml", nil)
	context.compile("javascript/exec.yaml", nil)
	context.compile("javascript/functions.yaml", nil)
	context.compile("openstack/hello-world.yaml", nil)
	context.compile("bpmn/open-loop.yaml", nil)
	context.compile("cloudify/advanced-blueprint-example.yaml", map[string]interface{}{
		"host_ip":                "1.2.3.4",
		"agent_user":             "my_user",
		"agent_private_key_path": "my_key",
	})
	context.compile("cloudify/example.yaml", nil)
	context.compile("hot/hello-world.yaml", map[string]interface{}{
		"username": "test",
	})
}

//
// Context
//

type Context struct {
	t             *testing.T
	root          string
	urlContext    *urlpkg.Context
	parserContext *parser.Context
}

func NewContext(t *testing.T) *Context {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("PUCCINI_TEST_ROOT"); !ok {
		root, _ = os.Getwd()
	}

	return &Context{
		t:             t,
		root:          root,
		urlContext:    urlpkg.NewContext(),
		parserContext: parser.NewContext(),
	}
}

func (self *Context) compile(url string, inputs map[string]interface{}) {
	self.t.Run(url, func(t *testing.T) {
		// Running the tests in parallel is not just for speed;
		// it actually helps us to find concurrency bugs
		t.Parallel()

		var serviceTemplate *normal.ServiceTemplate
		var clout *cloutpkg.Clout
		var problems *problemspkg.Problems
		var err error

		url_, err := urlpkg.NewURL(fmt.Sprintf("%s/examples/%s", self.root, url), self.urlContext)
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}

		if _, serviceTemplate, problems, err = self.parserContext.Parse(url_, nil, nil, inputs); err != nil {
			t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
			return
		}

		if clout, err = serviceTemplate.Compile(true); err != nil {
			t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
			return
		}

		js.Resolve(clout, problems, self.urlContext, true, "yaml", false, false, true)
		if !problems.Empty() {
			t.Errorf("%s", problems.ToString(true))
			return
		}

		js.Coerce(clout, problems, self.urlContext, true, "yaml", false, false, true)
		if !problems.Empty() {
			t.Errorf("%s", problems.ToString(true))
			return
		}
	})
}
