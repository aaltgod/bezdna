package domain

type Flag struct {
	ID        int64
	StreamID  int64
	Text      string
	Direction FlagDirection
}

type FlagRegexp struct {
	Regexp      string
	ServiceName string
}

type FlagDirection string

const (
	FlagDirection_UNKNOWN = "UNKNOWN"
	FlagDirection_IN      = "IN"
	FlagDirection_OUT     = "OUT"
)

func (fa FlagDirection) String() string {
	switch fa {
	case FlagDirection_IN, FlagDirection_OUT:
		return string(fa)
	default:
		return string(FlagDirection_UNKNOWN)
	}
}
