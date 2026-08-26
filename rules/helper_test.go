package rules_test

import "errors"

func errorf(msg string) error { return errors.New(msg) }
