package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceCreateOrUpdateParameters interface {
}

// RawEventSourceCreateOrUpdateParametersImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventSourceCreateOrUpdateParametersImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalEventSourceCreateOrUpdateParametersImplementation(input []byte) (EventSourceCreateOrUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceCreateOrUpdateParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	out := RawEventSourceCreateOrUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
