package valueobjects

func newNonEmptyString(value string, errIfEmpty error) (string, error) {
	if value == "" {
		return "", errIfEmpty
	}
	return value, nil
}
