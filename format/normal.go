package format

func Normalize(data interface{}) (interface{}, error) {
	code, err := EncodeYaml(data)
	if err != nil {
		return nil, err
	}
	return DecodeYaml(code)
}
