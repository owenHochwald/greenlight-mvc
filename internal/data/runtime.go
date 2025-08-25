package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

var invalidRuntimeFormat = errors.New("Invalid runtime format")

func (r *Runtime) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))

	if err != nil {
		return invalidRuntimeFormat
	}

	// split into parts
	parts := strings.Split(unquoted, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return invalidRuntimeFormat
	}

	// parse minutes
	i, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return invalidRuntimeFormat
	}
	*r = Runtime(i)
	return nil
}

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotes := strconv.Quote(jsonValue)

	return []byte(quotes), nil
}
