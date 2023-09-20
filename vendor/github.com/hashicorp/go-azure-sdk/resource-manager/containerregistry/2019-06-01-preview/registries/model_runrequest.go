package registries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunRequest interface {
}

// RawRunRequestImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRunRequestImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalRunRequestImplementation(input []byte) (RunRequest, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RunRequest into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "DockerBuildRequest") {
		var out DockerBuildRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DockerBuildRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EncodedTaskRunRequest") {
		var out EncodedTaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EncodedTaskRunRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileTaskRunRequest") {
		var out FileTaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileTaskRunRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TaskRunRequest") {
		var out TaskRunRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TaskRunRequest: %+v", err)
		}
		return out, nil
	}

	out := RawRunRequestImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
