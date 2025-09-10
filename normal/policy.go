package normal

//
// Policy
//

type Policy struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`
	Name            string           `json:"-" yaml:"-"`

	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
	Description string            `json:"description" yaml:"description"`
	Types       EntityTypes       `json:"types" yaml:"types"`
	Properties  Values            `json:"properties" yaml:"properties"`

	GroupTargets        []*Group         `json:"-" yaml:"-"`
	NodeTemplateTargets []*NodeTemplate  `json:"-" yaml:"-"`
	Triggers            []*PolicyTrigger `json:"-" yaml:"-"`
}

func (self *ServiceTemplate) NewPolicy(name string) *Policy {
	policy := &Policy{
		ServiceTemplate:     self,
		Name:                name,
		Metadata:            make(map[string]string),
		Types:               make(EntityTypes),
		Properties:          make(Values),
		GroupTargets:        make([]*Group, 0),
		NodeTemplateTargets: make([]*NodeTemplate, 0),
		Triggers:            make([]*PolicyTrigger, 0),
	}
	self.Policies[name] = policy
	return policy
}

//
// Policies
//

type Policies map[string]*Policy

//
// PolicyTrigger
//

type PolicyTrigger struct {
	Policy *Policy `json:"-" yaml:"-"`

	// TOSCA 2.0 specification fields
	Description string        `json:"description" yaml:"description"`
	Event       string        `json:"event" yaml:"event"`
	Condition   *FunctionCall `json:"condition" yaml:"condition"`

	// Legacy fields for backward compatibility
	EventType string     `json:"eventType" yaml:"eventType"`
	Operation *Operation `json:"operation" yaml:"operation"`
	Workflow  *Workflow  `json:"workflow" yaml:"workflow"`
}

func (self *Policy) NewTrigger() *PolicyTrigger {
	trigger := &PolicyTrigger{
		Policy: self,
	}
	self.Triggers = append(self.Triggers, trigger)
	return trigger
}
