package normal

//
// Notification
//

type Notification struct {
	Interface *Interface `json:"-" yaml:"-"`
	Name      string     `json:"-" yaml:"-"`

	Description    string   `json:"description" yaml:"description"`
	Implementation string   `json:"implementation" yaml:"implementation"`
	Dependencies   []string `json:"dependencies" yaml:"dependencies"`
	Timeout        int64    `json:"timeout" yaml:"timeout"`
	Host           string   `json:"host,omitempty" yaml:"host,omitempty"`
	Outputs        Mappings `json:"outputs" yaml:"outputs"`
}

func (self *Interface) NewNotification(name string) *Notification {
	notification := &Notification{
		Interface:    self,
		Name:         name,
		Dependencies: make([]string, 0),
		Outputs:      make(Mappings),
	}
	self.Notifications[name] = notification
	return notification
}

//
// Notifications
//

type Notifications map[string]*Notification
