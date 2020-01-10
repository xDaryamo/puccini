package format

func Normalize(data interface{}) (interface{}, error) {
	if code, err := EncodeYaml(data, " "); err == nil {
		return DecodeYaml(code)
	} else {
		return nil, err
	}
}
