package domain

import "time"

type Stream struct {
	Ack       uint64    `json:"ack"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}

type Fake struct {
	UserId    int64  `json:"userId"`
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
type GetStreamsByService struct {
	Service
	Offset int64 `json:"offset"`
	Limit  int16 `json:"limit"`
}
