package main

import "errors"

var (
	ErrValidation = errors.New("validation failed")
)

type stepErr struct {
	step string
	msg string
	cause error
}