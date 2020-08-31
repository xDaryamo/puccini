package profiles

import (
	bpmn_v1_0 "github.com/tliron/puccini/tosca/profiles/bpmn/v1_0"
	cloudify_v5_0_5 "github.com/tliron/puccini/tosca/profiles/cloudify/v5_0_5"
	common "github.com/tliron/puccini/tosca/profiles/common/v1_0"
	hot_v1_0 "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	implicit_v2_0 "github.com/tliron/puccini/tosca/profiles/implicit/v2_0"
	kubernetes_v1_0 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_0"
	openstack_v1_0 "github.com/tliron/puccini/tosca/profiles/openstack/v1_0"
	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_0 "github.com/tliron/puccini/tosca/profiles/simple/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
	simple_v1_2 "github.com/tliron/puccini/tosca/profiles/simple/v1_2"
	simple_v1_3 "github.com/tliron/puccini/tosca/profiles/simple/v1_3"

	"github.com/tliron/kutil/url"
)

func init() {
	initProfile(common.Profile)
	initProfile(simple_v1_0.Profile)
	initProfile(simple_v1_1.Profile)
	initProfile(simple_v1_2.Profile)
	initProfile(simple_v1_3.Profile)
	initProfile(implicit_v2_0.Profile)
	initProfile(simpleForNFV_v1_0.Profile)
	initProfile(kubernetes_v1_0.Profile)
	initProfile(openstack_v1_0.Profile)
	initProfile(bpmn_v1_0.Profile)
	initProfile(cloudify_v5_0_5.Profile)
	initProfile(hot_v1_0.Profile)
}

func initProfile(profile map[string]string) {
	for path, content := range profile {
		if err := url.RegisterInternalURL(path, content); err != nil {
			panic(err)
		}
	}
}
