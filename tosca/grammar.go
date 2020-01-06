package tosca

//
// GrammarVersion
//

type GrammarVersion struct {
	Version             string
	ImplicitProfilePath string
}

type GrammarVersions map[string][]GrammarVersion

func (self GrammarVersions) Add(keyword string, version string, implicitProfilePath string) {
	grammarVersion := GrammarVersion{version, implicitProfilePath}
	self[keyword] = append(self[keyword], grammarVersion)
}

//
// Grammar
//

type Grammar struct {
	Versions GrammarVersions
	Readers  Readers
}

func NewGrammar() Grammar {
	return Grammar{
		Versions: make(GrammarVersions),
		Readers:  make(Readers),
	}
}

func (self *Grammar) RegisterVersion(keyword string, version string, implicitProfilePath string) {
	self.Versions.Add(keyword, version, implicitProfilePath)
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
