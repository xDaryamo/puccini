package ard

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func KeyString(data interface{}) string {
	if s, ok := data.(string); ok {
		return s
	} else {
		return fmt.Sprintf("%s", data)
	}
}

//
// Key
//

type Key interface {
	GetKeyData() interface{}
}

//
// YamlKey
//

type YamlKey struct {
	Data interface{}
	Text string
}

func NewYamlKey(data interface{}) (*YamlKey, error) {
	var writer strings.Builder
	encoder := yaml.NewEncoder(&writer)
	if err := encoder.Encode(data); err == nil {
		text := writer.String()
		return &YamlKey{
			Data: data,
			Text: text,
		}, nil
	} else {
		return nil, err
	}
}

// Key interface
func (self *YamlKey) GetKeyData() interface{} {
	return self.Data
}

// fmt.Stringify interface
func (self *YamlKey) String() string {
	return self.Text
}
