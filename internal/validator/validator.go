package validator

import "regexp"

var EmailRX = regexp.MustCompile("\\w+@\\w+\\.\\w+")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}
func (v *Validator) AddError(key string, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}
func (v *Validator) String() string {
	str := ""
	for key, message := range v.Errors {
		str += key + ": " + message + "\n"
	}
	return str
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
