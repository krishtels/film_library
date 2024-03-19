package tools

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	violations []string
}

func (ve *ValidationError) NoViolations() bool {
	return len(ve.violations) == 0
}

func (ve *ValidationError) AddViolation(v string) {
	ve.violations = append(ve.violations, v)
}

func (ve *ValidationError) Error() string {
	return strings.Join(ve.violations[:], "; ")
}

type QueryableObject struct {
	keys   []string
	values []any
}

func NewQueryableObject() *QueryableObject {
	return &QueryableObject{
		keys:   make([]string, 0),
		values: make([]any, 0),
	}
}

func (qo *QueryableObject) IsEmpty() bool {
	return len(qo.keys) == 0
}

func (qo *QueryableObject) Add(key string, val any) {
	qo.keys = append(qo.keys, key)
	qo.values = append(qo.values, val)
}

func (qo *QueryableObject) Args(start int) string {
	args := make([]string, len(qo.keys))
	for i, v := range qo.keys {
		args[i] = fmt.Sprintf("%s = $%d", v, i+start)
	}

	return strings.Join(args, ", ")
}

func (qo *QueryableObject) Values() []any {
	return qo.values
}

func (qo *QueryableObject) Len() int {
	return len(qo.keys)
}
