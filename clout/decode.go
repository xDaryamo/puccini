package clout

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/yamlkeys"
)

func Decode(map_ ard.Map) (*Clout, error) {
	clout := NewClout()

	if data, ok := map_["version"]; ok {
		if version, ok := data.(string); ok {
			clout.Version = version
		} else {
			return nil, fmt.Errorf("malformed clout: \"version\" not a string: %T", data)
		}
	} else {
		return nil, errors.New("malformed clout: no \"version\"")
	}

	if data, ok := map_["metadata"]; ok {
		if metadata, ok := data.(ard.Map); ok {
			clout.Metadata = ard.ToStringMap(metadata)
		} else {
			return nil, fmt.Errorf("malformed clout: \"metadata\" not a map: %T", data)
		}
	}

	if data, ok := map_["properties"]; ok {
		if properties, ok := data.(ard.Map); ok {
			clout.Properties = ard.ToStringMap(properties)
		} else {
			return nil, fmt.Errorf("malformed clout: \"properties\" not a map: %T", data)
		}
	}

	if data, ok := map_["vertexes"]; ok {
		if vertexes, ok := data.(ard.Map); ok {
			for key, data := range vertexes {
				if id, ok := key.(string); ok {
					if map_, ok := data.(ard.Map); ok {
						vertex := clout.NewVertex(id)

						if data, ok := map_["metadata"]; ok {
							if metadata, ok := data.(ard.Map); ok {
								vertex.Metadata = ard.ToStringMap(metadata)
							} else {
								return nil, fmt.Errorf("malformed vertex: \"metadata\" not a map: %T", data)
							}
						}

						if data, ok := map_["properties"]; ok {
							if properties, ok := data.(ard.Map); ok {
								vertex.Properties = ard.ToStringMap(properties)
							} else {
								return nil, fmt.Errorf("malformed vertex: \"properties\" not a map: %T", data)
							}
						}

						if data, ok := map_["edgesOut"]; ok {
							if edgesOut, ok := data.(ard.List); ok {
								for _, data := range edgesOut {
									if map_, ok := data.(ard.Map); ok {
										if data, ok := map_["targetID"]; ok {
											if targetId, ok := data.(string); ok {
												edge := vertex.NewEdgeToID(targetId)

												if data, ok := map_["metadata"]; ok {
													if metadata, ok := data.(ard.Map); ok {
														edge.Metadata = ard.ToStringMap(metadata)
													} else {
														return nil, fmt.Errorf("malformed edge: \"metadata\" not a map: %T", data)
													}
												}

												if data, ok := map_["properties"]; ok {
													if properties, ok := data.(ard.Map); ok {
														edge.Properties = ard.ToStringMap(properties)
													} else {
														return nil, fmt.Errorf("malformed edge: \"properties\" not a map: %T", data)
													}
												}
											}
										} else {
											return nil, errors.New("malformed edge: no \"targetID\"")
										}
									} else {
										return nil, fmt.Errorf("malformed edge: not a map: %T", data)
									}
								}
							} else {
								return nil, fmt.Errorf("malformed vertex: \"edgesOut\" not a list: %T", data)
							}
						}
					} else {
						return nil, fmt.Errorf("malformed vertex: not a map: %T", data)
					}
				} else {
					return nil, fmt.Errorf("malformed vertex: id not a string: %T", key)
				}
			}
		} else {
			return nil, fmt.Errorf("malformed clout: \"vertexes\" not a map: %T", data)
		}
	}

	return clout, nil
}

func DecodeJson(reader io.Reader) (*Clout, error) {
	var clout Clout
	var err error

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&clout); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return &clout, nil
}

func DecodeYaml(reader io.Reader) (*Clout, error) {
	var err error

	var map_ ard.Map
	map_, err = yamlkeys.Decode(reader)

	var clout *Clout
	if clout, err = Decode(map_); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return clout, nil
}

func DecodeXml(reader io.Reader) (*Clout, error) {
	var clout Clout
	var err error

	decoder := xml.NewDecoder(reader)

	if err = decoder.Decode(&clout); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return &clout, nil
}
