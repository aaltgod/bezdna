package handler

import "github.com/pkg/errors"

var (
	ErrMinOffset = errors.New("offset can't be < 1")
	ErrMaxOffset = errors.New("offset can't be > 20")
	ErrMaxLimit  = errors.New("limit can't be > 20")
)

const (
	WrapReadAll   = "io.ReadAll"
	WrapUnmarshal = "json.Unmarshal"
	WrapMarshal   = "json.Marshal"

	WrapCreateService       = "service.CreateService"
	WrapGetServices         = "service.GetServices"
	WrapGetStreamsByService = "service.GetStreamsByService"
)

func WrapfCreateService(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.CreateService")
}

func WrapfGetServices(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.GetServices")
}

func WrapfGetStreamsByService(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.GetStreamsByService")
}
