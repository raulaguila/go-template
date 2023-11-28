package helpers

import (
	"errors"
	"strconv"
	"time"
)

func DurationFromString(str string, factor time.Duration) (time.Duration, error) {
	converted, err := strconv.Atoi(str)
	if err != nil {
		return time.Duration(0), errors.New("invalid string")
	}

	return time.Duration(converted) * factor, nil
}
