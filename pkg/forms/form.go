package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// New to initialize a custom Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// EmailPattern Use the regexp.MustCompile function to parse a pattern and compile a
// regular expression for sanity checking the format of an email address.
var EmailPattern = regexp.MustCompile("^[a-zA-Z\\d.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z\\d](?:[a-zA-Z\\d-]{0,61}[a-zA-Z\\d])?(?:\\.[a-zA-Z\\d](?:[a-zA-Z\\d-]{0,61}[a-zA-Z\\d])?)*$")

// MatchesPattern method to check that a specific field in the form
// matches a regular expression.
func (f *Form) MatchesPattern(field string, patter *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !patter.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Required Implement a Required method to check that specific fields in the form // data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength Implement a MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the // appropriate message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// MinLength method to check that a specific field in the form contains a minimum number of characters.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characteres)", d))
	}
}

// PermittedValues Implement a PermittedValues method to check that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Valid Implement a Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
