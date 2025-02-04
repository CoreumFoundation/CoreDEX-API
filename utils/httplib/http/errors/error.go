package httperror

import "fmt"

type Error struct {
	Name    ErrorName
	Details any
}

func (e *Error) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s details: %v", e.Name, e.Details)
	}

	return string(e.Name)
}

type ErrorName string

func NewErrorWithDetails(name ErrorName, details any) *Error {
	return &Error{
		Name:    name,
		Details: details,
	}
}

func NewError(name ErrorName) *Error {
	return &Error{
		Name: name,
	}
}
