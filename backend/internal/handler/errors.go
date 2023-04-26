package handler

import "github.com/pkg/errors"

var (
	ErrMinOffset = errors.New("offset can't be < 1")
	ErrMaxOffset = errors.New("offset can't be > 20")
)

const (
	WrapReadAll   = "io.ReadAll"
	WrapUnmarshal = "json.Unmarshal"
	WrapMarshal   = "json.Marshal"

	WrapAddService          = "service.AddService"
	WrapGetServices         = "service.GetServices"
	WrapGetStreamsByService = "service.GetStreamsByService"
)

func WrapfAddService(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.AddService")
}

func WrapfGetServices(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.GetServices")
}

func WrapfGetStreamsByService(err error, message string) error {
	return errors.Wrap(errors.Wrap(err, message), "handler.GetStreamsByService")
}
