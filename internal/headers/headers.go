package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func Create() Headers {
	return make(map[string]string)
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

// parses one key-value pair (header) at a time
func (h *Headers) Parse(data []byte) (int, bool, error) {
	i := bytes.Index(data, []byte("\r\n"))
	if i == -1 {
		return 0, false, nil
	}
	if i == 0 {
		return 2, true, nil
	}

	fieldLine := data[:i]

	i = bytes.Index(fieldLine, []byte(":"))
	if i == -1 {
		return 0, false, fmt.Errorf("Malformed field-line: %q\n", fieldLine)
	}

	fieldName := fieldLine[:i]
	fieldValue := fieldLine[i+1:]
	fieldName = bytes.TrimLeft(fieldName, " \t")
	fieldValue = bytes.TrimSpace(fieldValue)

	// Validate field-name:

	// Validate length (must be at least 1 character)
	if len(fieldName) == 0 {
		return 0, false, errors.New("Field-name can't be empty!")
	}

	// Validate characters (must be: A-Z, a-z, 0-9, !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~)
	for _, c := range fieldName {
		if (c < '0' || c > '9') && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c != '!' && c != '#' && c != '$' && c != '%' && c != '&' && c != '\'' && c != '*' && c != '+' && c != '-' && c != '.' && c != '^' && c != '_' && c != '`' && c != '|' && c != '~') {
			return 0, false, fmt.Errorf("Malformed field-name: %q\n", fieldName)
		}
	}

	headers := *h
	fieldName = bytes.ToLower(fieldName)
	if headers[string(fieldName)] == "" {
		headers[string(fieldName)] = string(fieldValue)
	} else {
		headers[string(fieldName)] += ", " + string(fieldValue)
	}

	return len(fieldLine) + 2, false, nil
}
