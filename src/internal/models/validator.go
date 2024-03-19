package models

import (
	"time"

	"film-library/src/internal/tools"
)

var sexMap = map[string]struct{}{
	"":       {},
	"female": {},
	"male":   {},
}

func ValidateFormatActorInfo(ai *ActorInfo) *tools.ValidationError {
	ve := &tools.ValidationError{}

	_, ok := sexMap[ai.Sex]
	if !ok {
		ve.AddViolation("incorrect sex format (expected one of [male, female])")
	}

	_, err := time.Parse(time.DateOnly, ai.Birthday)
	if err != nil && len(ai.Birthday) != 0 {
		ve.AddViolation("incorrect date format (expected format: 2006-01-02)")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateEmptyActorInfo(ai *ActorInfo) *tools.ValidationError {
	ve := &tools.ValidationError{}

	if len(ai.Name) == 0 {
		ve.AddViolation("name empty")
	}

	if len(ai.Sex) == 0 {
		ve.AddViolation("sex empty (expected one of [male, female])")
	}

	if len(ai.Birthday) == 0 {
		ve.AddViolation("date empty (expected format: 2006-01-02)")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}
