package parser

import (
	bpmn_v1_0 "github.com/tliron/puccini/tosca/profiles/bpmn/v1_0"
	cloudify_v4_5 "github.com/tliron/puccini/tosca/profiles/cloudify/v4_5"
	common "github.com/tliron/puccini/tosca/profiles/common/v1_0"
	hot_v1_0 "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	kubernetes_v1_0 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_0"
	openstack_v1_0 "github.com/tliron/puccini/tosca/profiles/openstack/v1_0"
	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
	simple_v1_2 "github.com/tliron/puccini/tosca/profiles/simple/v1_2"
	simple_v1_3 "github.com/tliron/puccini/tosca/profiles/simple/v1_3"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/url"
)

var InternalProfilePaths = map[string]map[string]string{
	"tosca_definitions_version": {
		"tosca_simple_yaml_1_0":            simple_v1_1.ProfileInternalPath, // TODO: properly support 1.0
		"tosca_simple_yaml_1_1":            simple_v1_1.ProfileInternalPath,
		"tosca_simple_yaml_1_2":            simple_v1_2.ProfileInternalPath,
		"tosca_simple_yaml_1_3":            simple_v1_3.ProfileInternalPath, // TODO
		"tosca_simple_profile_for_nfv_1_0": simpleForNFV_v1_0.ProfileInternalPath,
		"cloudify_dsl_1_3":                 cloudify_v4_5.ProfileInternalPath,
	},
}

func init() {
	initProfile(common.Profile)
	initProfile(simple_v1_1.Profile)
	initProfile(simple_v1_2.Profile)
	initProfile(simple_v1_3.Profile)
	initProfile(simpleForNFV_v1_0.Profile)
	initProfile(kubernetes_v1_0.Profile)
	initProfile(openstack_v1_0.Profile)
	initProfile(bpmn_v1_0.Profile)
	initProfile(cloudify_v4_5.Profile)
	initProfile(hot_v1_0.Profile)
}

func initProfile(profile map[string]string) {
	for k, v := range profile {
		url.Internal[k] = v
	}
}

func GetProfileImportSpec(context *tosca.Context) (*tosca.ImportSpec, bool) {
	if versionContext, version := GetVersion(context); version != nil {
		if paths, ok := InternalProfilePaths[versionContext.Name]; ok {
			if path, ok := paths[*version]; ok {
				if url_, err := url.NewValidInternalURL(path); err == nil {
					return &tosca.ImportSpec{url_, nil, true}, true
				} else {
					context.ReportError(err)
				}
			}
		}
	}

	return nil, false
}
