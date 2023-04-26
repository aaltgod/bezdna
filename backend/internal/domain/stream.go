package domain

import "time"

type Stream struct {
	Ack       uint64    `json:"ack"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}

type GetStreamsByService struct {
	Service
	Offset int64 `json:"offset"`
	Limit  int16 `json:"limit"`
}
