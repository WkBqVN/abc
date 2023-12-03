package controller

import "errors"

var (
	ErrConnectDatabaseFailed = errors.New("failed to connect data base")
	ErrRequestDataInvailid   = errors.New("request data is incorrect")
)
