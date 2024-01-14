package service

import "errors"

const (
	WrapUpsertService       = "dbRepository.UpsertService"
	WrapGetServices         = "dbRepository.GetServices"
	WrapGetStreamsByService = "dbRepository.GetStreamsByService"

	WrapAddConfig = "sniffer.AddConfig"
)

var (
	ErrAlreadyExist = errors.New("already exist")
)
