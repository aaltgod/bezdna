package handler

import "github.com/pkg/errors"

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
