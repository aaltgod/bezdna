package domain

type Service struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	FlagRegexp string `json:"flag_regexp"`
}
