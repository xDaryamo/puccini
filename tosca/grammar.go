package tosca

//
// GrammarVersion
//

type GrammarVersion struct {
	Version             string
	ProfileInternalPath string
}

type GrammarVersions map[string][]GrammarVersion

func (self GrammarVersions) Add(keyword string, version string, profileInternalPath string) {
	grammarVersion := GrammarVersion{version, profileInternalPath}
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

func (self *Grammar) RegisterVersion(keyword string, version string, profileInternalPath string) {
	self.Versions.Add(keyword, version, profileInternalPath)
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
