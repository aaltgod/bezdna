package service

import "errors"

const (
	WrapInsertService       = "dbRepository.InsertService"
	WrapGetServices         = "dbRepository.GetServices"
	WrapGetStreamsByService = "dbRepository.GetStreamsByService"

	WrapAddConfig = "sniffer.AddConfig"
)

var (
	ErrAlreadyExist = errors.New("already exist")
)
