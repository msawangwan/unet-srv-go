package handler

import "errors"

var (
	errNilBody = errors.New("expected a body, recvd nil instead")
)
