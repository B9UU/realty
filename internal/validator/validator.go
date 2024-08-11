package validator

import (
	"regexp"
)

// var (
// 	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
// )

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// contains errors
type Validator struct {
	Errors map[string]string
}

// initiate new Validatore
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// check if there's no erros in Validator.Error
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// adds erros to Validator.Errors
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// if the comparison operation is false adds key, message to the validator.Errors map
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// check if value in list
func In(value string, list ...string) bool {
	for _, l := range list {
		if value == l {
			return true
		}
	}
	return false
}

// check if value matches the regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// check if the values are unique
func Unique(values []string) bool {
	unique := make(map[string]bool)
	for _, v := range values {
		if unique[v] {
			return false
		}
		unique[v] = true
	}
	return true
}
