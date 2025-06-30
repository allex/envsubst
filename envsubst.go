package envsubst

import (
	"io/ioutil"
	"os"

	"github.com/allex/envsubst/parse"
)

// String returns the parsed template string after processing it.
// If the parser encounters invalid input, it returns an error describing the failure.
func String(s string) (string, error) {
	return StringRestricted(s, false, false)
}

// StringRestricted returns the parsed template string after processing it.
// If the parser encounters invalid input, or a restriction is violated, it returns
// an error describing the failure.
// Errors on first failure or returns a collection of failures if failOnFirst is false
func StringRestricted(s string, noUnset, noEmpty bool) (string, error) {
	return StringRestrictedNoDigit(s, noUnset, noEmpty, false)
}

// Like StringRestricted but additionally allows to ignore env variables which start with a digit.
func StringRestrictedNoDigit(s string, noUnset, noEmpty bool, noDigit bool) (string, error) {
	return StringRestrictedKeepUnset(s, noUnset, noEmpty, noDigit, false)
}

// StringRestrictedKeepUnset provides full control over all restriction options including KeepUnset.
// If keepUnset is true, undefined variables will be kept as their original text instead of being substituted or causing errors.
func StringRestrictedKeepUnset(s string, noUnset, noEmpty bool, noDigit bool, keepUnset bool) (string, error) {
	return parse.New("string", parse.NewEnv(os.Environ()),
		&parse.Restrictions{NoUnset: noUnset, NoEmpty: noEmpty, NoDigit: noDigit, KeepUnset: keepUnset, VarMatcher: nil}).Parse(s)
}

// Bytes returns the bytes represented by the parsed template after processing it.
// If the parser encounters invalid input, it returns an error describing the failure.
func Bytes(b []byte) ([]byte, error) {
	return BytesRestricted(b, false, false)
}

// BytesRestricted returns the bytes represented by the parsed template after processing it.
// If the parser encounters invalid input, or a restriction is violated, it returns
// an error describing the failure.
func BytesRestricted(b []byte, noUnset, noEmpty bool) ([]byte, error) {
	return BytesRestrictedNoDigit(b, noUnset, noEmpty, false)
}

// Like BytesRestricted but additionally allows to ignore env variables which start with a digit.
func BytesRestrictedNoDigit(b []byte, noUnset, noEmpty bool, noDigit bool) ([]byte, error) {
	return BytesRestrictedKeepUnset(b, noUnset, noEmpty, noDigit, false)
}

// BytesRestrictedKeepUnset provides full control over all restriction options including KeepUnset.
// If keepUnset is true, undefined variables will be kept as their original text instead of being substituted or causing errors.
func BytesRestrictedKeepUnset(b []byte, noUnset, noEmpty bool, noDigit bool, keepUnset bool) ([]byte, error) {
	s, err := parse.New("bytes", parse.NewEnv(os.Environ()),
		&parse.Restrictions{NoUnset: noUnset, NoEmpty: noEmpty, NoDigit: noDigit, KeepUnset: keepUnset, VarMatcher: nil}).Parse(string(b))
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// ReadFile call io.ReadFile with the given file name.
// If the call to io.ReadFile failed it returns the error; otherwise it will
// call envsubst.Bytes with the returned content.
func ReadFile(filename string) ([]byte, error) {
	return ReadFileRestricted(filename, false, false)
}

// ReadFileRestricted calls io.ReadFile with the given file name.
// If the call to io.ReadFile failed it returns the error; otherwise it will
// call envsubst.Bytes with the returned content.
func ReadFileRestricted(filename string, noUnset, noEmpty bool) ([]byte, error) {
	return ReadFileRestrictedNoDigit(filename, noUnset, noEmpty, false)
}

// Like ReadFileRestricted but additionally allows to ignore env variables which start with a digit.
func ReadFileRestrictedNoDigit(filename string, noUnset, noEmpty bool, noDigit bool) ([]byte, error) {
	return ReadFileRestrictedKeepUnset(filename, noUnset, noEmpty, noDigit, false)
}

// ReadFileRestrictedKeepUnset provides full control over all restriction options including KeepUnset.
// If keepUnset is true, undefined variables will be kept as their original text instead of being substituted or causing errors.
func ReadFileRestrictedKeepUnset(filename string, noUnset, noEmpty bool, noDigit bool, keepUnset bool) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return BytesRestrictedKeepUnset(b, noUnset, noEmpty, noDigit, keepUnset)
}
