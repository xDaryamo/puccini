package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Plugin
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-plugins/]
//

type Plugin struct {
	*Entity `name:"plugin"`
	Name    string `namespace:""`

	Executor            *string `read:"executor" require:""`
	Source              *string `read:"source"`
	InstallArguments    *string `read:"install_arguments"`
	Install             *bool   `read:"install"`
	PackageName         *string `read:"package_name"`
	PackageVersion      *string `read:"package_version"`
	SupportedPlatform   *string `read:"supported_platform"`
	Distribution        *string `read:"distribution"`
	DistributionVersion *string `read:"distribution_version"`
	DistributionRelease *string `read:"distribution_release"`
}

func NewPlugin(context *tosca.Context) *Plugin {
	return &Plugin{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadPlugin(context *tosca.Context) tosca.EntityPtr {
	self := NewPlugin(context)

	context.ValidateUnsupportedFields(context.ReadFields(self))

	if self.Install == nil {
		// Default to true
		install := true
		self.Install = &install
	}

	if self.Executor != nil {
		executor := *self.Executor
		switch executor {
		case "central_deployment_agent", "host_agent":
		default:
			context.FieldChild("executor", executor).ReportFieldUnsupportedValue()
		}
	}

	if *self.Install && (self.Source == nil) && (self.PackageName == nil) {
		context.FieldChild("source", nil).ReportFieldMissing()
		context.FieldChild("package_name", nil).ReportFieldMissing()
	}

	return self
}

//
// Plugins
//

type Plugins []*Plugin
