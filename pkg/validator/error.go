package validator

import "fmt"

var ErrValidator *ValidatorError

type ValidatorError struct {
	field string
	tag   string
	param string
	value interface{}
}

func (m *ValidatorError) Error() string {
	return fmt.Sprintf("%s does not meet the '%s[%s]' requirement with value '%v'", m.field, m.tag, m.param, m.value)
}
