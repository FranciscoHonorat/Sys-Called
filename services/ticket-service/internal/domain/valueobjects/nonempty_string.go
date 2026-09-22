package valueobjects

import "encoding/json"

func newNonEmptyString(value string, errIfEmpty error) (string, error) {
	if value == "" {
		return "", errIfEmpty
	}
	return value, nil
}

func marshalNonEmptyString(value string) ([]byte, error) {
	return json.Marshal(value)
}

func unmarshalNonEmptyString(data []byte, errIfEmpty error) (string, error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}
	return newNonEmptyString(value, errIfEmpty)
}
