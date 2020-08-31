package hot

import (
	"strings"

	"github.com/tliron/kutil/ard"
)

//
// Parameter types
//

var ParameterTypes = []string{
	"boolean",
	"comma_delimited_list",
	"json",
	"number",
	"string",
}

func IsParameterTypeValid(type_ string) bool {
	for _, t := range ParameterTypes {
		if t == type_ {
			return true
		}
	}
	return false
}

func (self *Value) ValidateParameterType(type_ string) bool {
	switch type_ {
	case "boolean":
		return self.Context.ValidateType(ard.TypeBoolean)
	case "comma_delimited_list":
		if self.Context.ValidateType(ard.TypeList) {
			for index, e := range self.Context.Data.(ard.List) {
				if !self.Context.ListChild(index, e).ValidateType(ard.TypeString) {
					return false
				}
			}
			return true
		} else {
			return false
		}
	case "json":
		return self.Context.ValidateType(ard.TypeMap, ard.TypeList)
	case "number":
		return self.Context.ValidateType(ard.TypeInteger, ard.TypeFloat)
	case "string":
		return self.Context.ValidateType(ard.TypeString)
	default:
		panic("unsupported parameter type")
	}
}

func (self *Value) CoerceParameterType(type_ string) {
	switch type_ {
	case "boolean":
		switch data := self.Context.Data.(type) {
		case string:
			switch data {
			case "t", "true", "on", "y", "yes", "1":
				self.Context.Data = true
			case "f", "false", "off", "n", "no", "0":
				self.Context.Data = false
			}

		case int:
			switch data {
			case 1:
				self.Context.Data = true
			case 0:
				self.Context.Data = false
			}
		}

	case "comma_delimited_list":
		switch data := self.Context.Data.(type) {
		case string:
			split := strings.Split(data, ",")
			list := make(ard.List, len(split))
			for index, s := range split {
				list[index] = s
			}
			self.Context.Data = list
		}
	}
}

//
// Resource types
//
// [https://docs.openstack.org/heat/stein/template_guide/openstack.html]
//

var ResourceTypes = map[string]string{
	// Aodh
	"OS::Aodh::CompositeAlarm":                     "openstack:aodh.CompositeAlarm",
	"OS::Aodh::EventAlarm":                         "openstack:aodh.EventAlarm",
	"OS::Aodh::GnocchiAggregationByMetricsAlarm":   "openstack:aodh.GnocchiAggregationByMetricsAlarm",
	"OS::Aodh::GnocchiAggregationByResourcesAlarm": "openstack:aodh.GnocchiAggregationByResourcesAlarm",
	"OS::Aodh::GnocchiResourcesAlarm":              "openstack:aodh.GnocchiResourcesAlarm",
	// Barbican
	"OS::Barbican::CertificateContainer": "openstack:barbican.CertificateContainer",
	"OS::Barbican::GenericContainer":     "openstack:barbican.GenericContainer",
	"OS::Barbican::Order":                "openstack:barbican.Order",
	"OS::Barbican::RSAContainer":         "openstack:barbican.RSAContainer",
	"OS::Barbican::Secret":               "openstack:barbican.Secret",
	// Cinder
	"OS::Cinder::EncryptedVolumeType": "openstack:cinder.EncryptedVolumeType",
	"OS::Cinder::QoSAssociation":      "openstack:cinder.QoSAssociation",
	"OS::Cinder::QoSSpecs":            "openstack:cinder.QoSSpecs",
	"OS::Cinder::Quota":               "openstack:cinder.Quota",
	"OS::Cinder::Volume":              "openstack:cinder.Volume",
	"OS::Cinder::VolumeAttachment":    "openstack:cinder.VolumeAttachment",
	"OS::Cinder::VolumeType":          "openstack:cinder.VolumeType",
	// Designate
	"OS::Designate::RecordSet": "openstack:designate.RecordSet",
	"OS::Designate::Zone":      "openstack:designate.Zone",
	// Heat
	"OS::Heat::AccessPolicy":              "openstack:heat.AccessPolicy",
	"OS::Heat::AutoScalingGroup":          "openstack:heat.AutoScalingGroup",
	"OS::Heat::CloudConfig":               "openstack:heat.CloudConfig",
	"OS::Heat::Delay":                     "openstack:heat.Delay",
	"OS::Heat::DeployedServer":            "openstack:heat.DeployedServer",
	"OS::Heat::InstanceGroup":             "openstack:heat.InstanceGroup",
	"OS::Heat::MultipartMime":             "openstack:heat.MultipartMime",
	"OS::Heat::None":                      "openstack:heat.None",
	"OS::Heat::RandomString":              "openstack:heat.RandomString",
	"OS::Heat::ResourceChain":             "openstack:heat.ResourceChain",
	"OS::Heat::ResourceGroup":             "openstack:heat.ResourceGroup",
	"OS::Heat::ScalingPolicy":             "openstack:heat.ScalingPolicy",
	"OS::Heat::SoftwareComponent":         "openstack:heat.SoftwareComponent",
	"OS::Heat::SoftwareConfig":            "openstack:heat.SoftwareConfig",
	"OS::Heat::SoftwareDeployment":        "openstack:heat.SoftwareDeployment",
	"OS::Heat::SoftwareDeploymentGroup":   "openstack:heat.SoftwareDeploymentGroup",
	"OS::Heat::Stack":                     "openstack:heat.Stack",
	"OS::Heat::StructuredConfig":          "openstack:heat.StructuredConfig",
	"OS::Heat::StructuredDeployment":      "openstack:heat.StructuredDeployment",
	"OS::Heat::StructuredDeploymentGroup": "openstack:heat.StructuredDeploymentGroup",
	"OS::Heat::SwiftSignal":               "openstack:heat.SwiftSignal",
	"OS::Heat::SwiftSignalHandle":         "openstack:heat.SwiftSignalHandle",
	"OS::Heat::TestResource":              "openstack:heat.TestResource",
	"OS::Heat::UpdateWaitConditionHandle": "openstack:heat.UpdateWaitConditionHandle",
	"OS::Heat::Value":                     "openstack:heat.Value",
	"OS::Heat::WaitCondition":             "openstack:heat.WaitCondition",
	"OS::Heat::WaitConditionHandle":       "openstack:heat.WaitConditionHandle",
	// Keystone
	"OS::Keystone::Domain":              "openstack:keystone.Domain",
	"OS::Keystone::Endpoint":            "openstack:keystone.Endpoint",
	"OS::Keystone::Group":               "openstack:keystone.Group",
	"OS::Keystone::GroupRoleAssignment": "openstack:keystone.GroupRoleAssignment",
	"OS::Keystone::Project":             "openstack:keystone.Project",
	"OS::Keystone::Region":              "openstack:keystone.Region",
	"OS::Keystone::Role":                "openstack:keystone.Role",
	"OS::Keystone::Service":             "openstack:keystone.Service",
	"OS::Keystone::User":                "openstack:keystone.User",
	"OS::Keystone::UserRoleAssignment":  "openstack:keystone.UserRoleAssignment",
	// Magnum
	"OS::Magnum::Cluster":         "openstack:magnum.Cluster",
	"OS::Magnum::ClusterTemplate": "openstack:magnum.ClusterTemplate",
	// Manila
	"OS::Manila::SecurityService": "openstack:manila.SecurityService",
	"OS::Manila::Share":           "openstack:manila.Share",
	"OS::Manila::ShareNetwork":    "openstack:manila.ShareNetwork",
	"OS::Manila::ShareType":       "openstack:manila.ShareType",
	// Mistral
	"OS::Mistral::CronTrigger":      "openstack:mistral.CronTrigger",
	"OS::Mistral::ExternalResource": "openstack:mistral.ExternalResource",
	"OS::Mistral::Workflow":         "openstack:mistral.Workflow",
	// Monasca
	"OS::Monasca::AlarmDefinition": "openstack:monasca.AlarmDefinition",
	"OS::Monasca::Notification":    "openstack:monasca.Notification",
	// Neutron
	"OS::Neutron::AddressScope":          "openstack:neutron.AddressScope",
	"OS::Neutron::Firewall":              "openstack:neutron.Firewall",
	"OS::Neutron::FirewallPolicy":        "openstack:neutron.FirewallPolicy",
	"OS::Neutron::FirewallRule":          "openstack:neutron.FirewallRule",
	"OS::Neutron::FloatingIP":            "openstack:neutron.FloatingIP",
	"OS::Neutron::FloatingIPAssociation": "openstack:neutron.FloatingIPAssociation",
	"OS::Neutron::IKEPolicy":             "openstack:neutron.IKEPolicy",
	"OS::Neutron::IPsecPolicy":           "openstack:neutron.IPsecPolicy",
	"OS::Neutron::IPsecSiteConnection":   "openstack:neutron.IPsecSiteConnection",
	"OS::Neutron::LBaaS::HealthMonitor":  "openstack:neutron.lbaas.HealthMonitor",
	"OS::Neutron::LBaaS::L7Policy":       "openstack:neutron.lbaas.L7Policy",
	"OS::Neutron::LBaaS::L7Rule":         "openstack:neutron.lbaas.L7Rule",
	"OS::Neutron::LBaaS::Listener":       "openstack:neutron.lbaas.Listener",
	"OS::Neutron::LBaaS::LoadBalancer":   "openstack:neutron.lbaas.LoadBalancer",
	"OS::Neutron::LBaaS::Pool":           "openstack:neutron.lbaas.Pool",
	"OS::Neutron::LBaaS::PoolMember":     "openstack:neutron.lbaas.PoolMember",
	"OS::Neutron::MeteringLabel":         "openstack:neutron.MeteringLabel",
	"OS::Neutron::MeteringRule":          "openstack:neutron.MeteringRule",
	"OS::Neutron::Net":                   "openstack:neutron.Net",
	"OS::Neutron::NetworkGateway":        "openstack:neutron.NetworkGateway",
	"OS::Neutron::Port":                  "openstack:neutron.Port",
	"OS::Neutron::ProviderNet":           "openstack:neutron.ProviderNet",
	"OS::Neutron::QoSBandwidthLimitRule": "openstack:neutron.QoSBandwidthLimitRule",
	"OS::Neutron::QoSDscpMarkingRule":    "openstack:neutron.QoSDscpMarkingRule",
	"OS::Neutron::QoSPolicy":             "openstack:neutron.QoSPolicy",
	"OS::Neutron::Quota":                 "openstack:neutron.Quota",
	"OS::Neutron::RBACPolicy":            "openstack:neutron.RBACPolicy",
	"OS::Neutron::Router":                "openstack:neutron.Router",
	"OS::Neutron::RouterInterface":       "openstack:neutron.RouterInterface",
	"OS::Neutron::SecurityGroup":         "openstack:neutron.SecurityGroup",
	"OS::Neutron::SecurityGroupRule":     "openstack:neutron.SecurityGroupRuleSecurityGroupRule",
	"OS::Neutron::Segment":               "openstack:neutron.Segment",
	"OS::Neutron::Subnet":                "openstack:neutron.Subnet",
	"OS::Neutron::SubnetPool":            "openstack:neutron.SubnetPool",
	"OS::Neutron::Trunk":                 "openstack:neutron.Trunk",
	"OS::Neutron::VPNService":            "openstack:neutron.VPNService",
	// Nova
	"OS::Nova::Flavor":        "openstack:nova.Flavor",
	"OS::Nova::HostAggregate": "openstack:nova.HostAggregate",
	"OS::Nova::KeyPair":       "openstack:nova.KeyPair",
	"OS::Nova::Quota":         "openstack:nova.Quota",
	"OS::Nova::Server":        "openstack:nova.Server",
	"OS::Nova::ServerGroup":   "openstack:nova.ServerGroup",
	// Octavia
	"OS::Octavia::HealthMonitor": "openstack:octavia.HealthMonitor",
	"OS::Octavia::L7Policy":      "openstack:octavia.L7Policy",
	"OS::Octavia::L7Rule":        "openstack:octavia.L7Rule",
	"OS::Octavia::Listener":      "openstack:octavia.Listener",
	"OS::Octavia::LoadBalancer":  "openstack:octavia.LoadBalancer",
	"OS::Octavia::Pool":          "openstack:octavia.Pool",
	"OS::Octavia::PoolMember":    "openstack:octavia.PoolMember",
	// Sahara
	"OS::Sahara::Cluster":           "openstack:sahara.Cluster",
	"OS::Sahara::ClusterTemplate":   "openstack:sahara.ClusterTemplate",
	"OS::Sahara::DataSource":        "openstack:sahara.DataSource",
	"OS::Sahara::ImageRegistry":     "openstack:sahara.ImageRegistry",
	"OS::Sahara::Job":               "openstack:sahara.Job",
	"OS::Sahara::JobBinary":         "openstack:sahara.JobBinary",
	"OS::Sahara::NodeGroupTemplate": "openstack:sahara.NodeGroupTemplate",
	// Senlin
	"OS::Senlin::Cluster":  "openstack:senlin.Cluster",
	"OS::Senlin::Node":     "openstack:senlin.Node",
	"OS::Senlin::Policy":   "openstack:senlin.Policy",
	"OS::Senlin::Profile":  "openstack:senlin.Profile",
	"OS::Senlin::Receiver": "openstack:senlin.Receiver",
	// Swift
	"OS::Swift::Container": "openstack:swift.Container",
	// Trove
	"OS::Trove::Cluster":  "openstack:trove.Cluster",
	"OS::Trove::Instance": "openstack:trove.Instance",
	// Zaqar
	"OS::Zaqar::MistralTrigger": "openstack:zaqar.MistralTrigger",
	"OS::Zaqar::Queue":          "openstack:zaqar.Queue",
	"OS::Zaqar::SignedQueueURL": "openstack:zaqar.SignedQueueURL",
	"OS::Zaqar::Subscription":   "openstack:zaqar.Subscription",
	// Zun
	"OS::Zun::Container": "openstack:zun.Container",
}
